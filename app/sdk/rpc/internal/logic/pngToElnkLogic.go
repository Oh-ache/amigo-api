package logic

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"net/http"
	"strings"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	osContext "amigo-api/common/utils/plug/objectsave/context"
	"amigo-api/common/utils/plug/objectsave/factory"
	"amigo-api/common/utils/plug/objectsave/model"

	"github.com/nfnt/resize"
	"github.com/zeromicro/go-zero/core/logx"
)

type PngToElnkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPngToElnkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PngToElnkLogic {
	return &PngToElnkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PngToElnkLogic) PngToElnk(in *pb.PngToElnkReq) (*pb.PngToElnkResp, error) {
	// 从 URL 下载图片
	resp, err := http.Get(in.Url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解码 PNG 图片
	img, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 调整图片尺寸为 240x416
	resizedImg := resize.Resize(240, 416, img, resize.Lanczos3)

	// 预处理：gamma 校正 + 对比度拉伸 + 饱和度增强
	enhancedImg := enhanceForElnk(resizedImg)

	// 使用 Floyd-Steinberg 抖动算法转换为 4 色（黑、白、红、黄）格式
	convertedImg := floydSteinbergDithering(enhancedImg)

	// 编码为 BMP 格式
	bmpData, err := encodeToBMP(convertedImg)
	if err != nil {
		return nil, err
	}

	// 使用与 UploadFile 相同的 OSS 服务上传转换后的图片
	factory := factory.NewStorageFactory()

	storageConfig := &model.StorageConfig{
		Type: "oss",
		OssConfig: &model.OssConfig{
			Endpoint:        GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.endpoint"),
			AccessKeyId:     GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.accessKey"),
			AccessKeySecret: GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.accessKeySecret"),
			Bucket:          GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.bucket"),
			Region:          GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.region"),
		},
	}

	ossClient, _ := factory.CreateClient(storageConfig)
	storageCtx := osContext.NewStorageContext(ossClient)

	// 生成文件名（基于原始 URL 文件名）
	parts := strings.Split(in.Url, "/")
	filename := parts[len(parts)-1]
	if strings.HasSuffix(filename, ".png") {
		filename = strings.TrimSuffix(filename, ".png") + ".bmp"
	} else {
		filename += ".bmp"
	}

	_, err = storageCtx.UploadFile(filename, bmpData)
	if err != nil {
		return nil, err
	}

	host := GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.oss.host")
	uploadedUrl := fmt.Sprintf("%s/%s", host, filename)

	return &pb.PngToElnkResp{
		Url: uploadedUrl,
	}, nil
}

// Floyd-Steinberg 抖动算法实现
func floydSteinbergDithering(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	dst := image.NewRGBA(bounds)

	// 4 种墨水屏支持的颜色
	// 必须与 encodeToBMP 中的调色板保持一致
	colors := []color.RGBA{
		{0, 0, 0, 255},       // 黑色
		{255, 255, 255, 255}, // 白色
		{255, 0, 0, 255},     // 红色
		{255, 255, 0, 255},   // 黄色
	}

	// 使用 int16 误差缓冲区，避免浮点精度损失
	// 多一行缓冲区用于下一行的误差扩散
	errBuf := make([][3]int16, width*(height+1))
	for i := range errBuf {
		errBuf[i] = [3]int16{0, 0, 0}
	}

	idx := func(x, y int) int {
		return y*width + x
	}

	// 转换为 RGBA 图像并加上初始误差缓冲
	originalRGBA := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalRGBA.SetRGBA(x, y, color.RGBAModel.Convert(img.At(x, y)).(color.RGBA))
		}
	}

	// 应用 Floyd-Steinberg 抖动（蛇形扫描）
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if y%2 == 0 {
			// 偶数行：从左到右
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := originalRGBA.RGBAAt(x, y)
				ix := x - bounds.Min.X
				iy := y - bounds.Min.Y

				// 当前像素值加上累积误差
				r := clampI16(int16(c.R) + errBuf[idx(ix, iy)][0])
				g := clampI16(int16(c.G) + errBuf[idx(ix, iy)][1])
				b := clampI16(int16(c.B) + errBuf[idx(ix, iy)][2])

				// 找到最接近的颜色
				var bestColor color.RGBA
				var minDistance int = 0x7FFFFFFF
				srcColor := color.RGBA{uint8(r), uint8(g), uint8(b), c.A}
				for _, col := range colors {
					distance := colorDistance(srcColor, col)
					if distance < minDistance {
						minDistance = distance
						bestColor = col
					}
				}

				dst.SetRGBA(x, y, bestColor)

				// 计算误差
				errR := int16(int(srcColor.R) - int(bestColor.R))
				errG := int16(int(srcColor.G) - int(bestColor.G))
				errB := int16(int(srcColor.B) - int(bestColor.B))

				// Floyd-Steinberg 误差分布（使用整数运算，分母 16）
				//      X   7
				// 3   5   1
				if ix+1 < width {
					errBuf[idx(ix+1, iy)][0] += errR * 7 / 16
					errBuf[idx(ix+1, iy)][1] += errG * 7 / 16
					errBuf[idx(ix+1, iy)][2] += errB * 7 / 16
				}
				if iy+1 < height {
					if ix-1 >= 0 {
						errBuf[idx(ix-1, iy+1)][0] += errR * 3 / 16
						errBuf[idx(ix-1, iy+1)][1] += errG * 3 / 16
						errBuf[idx(ix-1, iy+1)][2] += errB * 3 / 16
					}
					errBuf[idx(ix, iy+1)][0] += errR * 5 / 16
					errBuf[idx(ix, iy+1)][1] += errG * 5 / 16
					errBuf[idx(ix, iy+1)][2] += errB * 5 / 16
					if ix+1 < width {
						errBuf[idx(ix+1, iy+1)][0] += errR * 1 / 16
						errBuf[idx(ix+1, iy+1)][1] += errG * 1 / 16
						errBuf[idx(ix+1, iy+1)][2] += errB * 1 / 16
					}
				}
			}
		} else {
			// 奇数行：从右到左（蛇形扫描）
			for x := bounds.Max.X - 1; x >= bounds.Min.X; x-- {
				c := originalRGBA.RGBAAt(x, y)
				ix := x - bounds.Min.X
				iy := y - bounds.Min.Y

				r := clampI16(int16(c.R) + errBuf[idx(ix, iy)][0])
				g := clampI16(int16(c.G) + errBuf[idx(ix, iy)][1])
				b := clampI16(int16(c.B) + errBuf[idx(ix, iy)][2])

				var bestColor color.RGBA
				var minDistance int = 0x7FFFFFFF
				srcColor := color.RGBA{uint8(r), uint8(g), uint8(b), c.A}
				for _, col := range colors {
					distance := colorDistance(srcColor, col)
					if distance < minDistance {
						minDistance = distance
						bestColor = col
					}
				}

				dst.SetRGBA(x, y, bestColor)

				errR := int16(int(srcColor.R) - int(bestColor.R))
				errG := int16(int(srcColor.G) - int(bestColor.G))
				errB := int16(int(srcColor.B) - int(bestColor.B))

				// 反向扫描的误差分布
				if ix-1 >= 0 {
					errBuf[idx(ix-1, iy)][0] += errR * 7 / 16
					errBuf[idx(ix-1, iy)][1] += errG * 7 / 16
					errBuf[idx(ix-1, iy)][2] += errB * 7 / 16
				}
				if iy+1 < height {
					if ix+1 < width {
						errBuf[idx(ix+1, iy+1)][0] += errR * 3 / 16
						errBuf[idx(ix+1, iy+1)][1] += errG * 3 / 16
						errBuf[idx(ix+1, iy+1)][2] += errB * 3 / 16
					}
					errBuf[idx(ix, iy+1)][0] += errR * 5 / 16
					errBuf[idx(ix, iy+1)][1] += errG * 5 / 16
					errBuf[idx(ix, iy+1)][2] += errB * 5 / 16
					if ix-1 >= 0 {
						errBuf[idx(ix-1, iy+1)][0] += errR * 1 / 16
						errBuf[idx(ix-1, iy+1)][1] += errG * 1 / 16
						errBuf[idx(ix-1, iy+1)][2] += errB * 1 / 16
					}
				}
			}
		}
	}

	return dst
}

// clampI16 将 int16 钳制到 0-255 范围
func clampI16(v int16) uint8 {
	if v < 0 {
		return 0
	} else if v > 255 {
		return 255
	}
	return uint8(v)
}

func clamp(v int) int {
	if v < 0 {
		return 0
	} else if v > 255 {
		return 255
	}
	return v
}

// RGB 转 YCbCr
func rgbToYCbCr(c color.RGBA) (y, cb, cr float64) {
	r := float64(c.R)
	g := float64(c.G)
	b := float64(c.B)
	y = 0.299*r + 0.587*g + 0.114*b
	cb = 128.0 - 0.168736*r - 0.331264*g + 0.5*b
	cr = 128.0 + 0.5*r - 0.418688*g - 0.081312*b
	return y, cb, cr
}

// 感知颜色距离（针对墨水屏 4 色场景优化）
// 墨水屏的黑/白靠亮度区分，红/黄靠色度区分
func colorDistance(c1, c2 color.RGBA) int {
	y1, cb1, cr1 := rgbToYCbCr(c1)
	y2, cb2, cr2 := rgbToYCbCr(c2)

	yDiff := y1 - y2
	cbDiff := cb1 - cb2
	crDiff := cr1 - cr2

	// 增强色度权重，避免红色被误判为黄色或黑色
	// 亮度权重 3x，色度权重各 2x
	return int(3*yDiff*yDiff + 2*cbDiff*cbDiff + 2*crDiff*crDiff)
}

// 图像预处理：针对墨水屏优化
// 1. Gamma 校正（提升暗部细节）
// 2. 对比度拉伸
// 3. 饱和度增强（帮助红/黄区分）
func enhanceForElnk(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// 第一步：计算灰度的最小和最大值
	var minGray, maxGray float64 = 255, 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			gray := float64(c.R)*0.299 + float64(c.G)*0.587 + float64(c.B)*0.114
			if gray < minGray {
				minGray = gray
			}
			if gray > maxGray {
				maxGray = gray
			}
		}
	}

	// 确保有足够的对比度范围
	if maxGray-minGray < 10 {
		minGray = 0
		maxGray = 255
	}

	// 第二步：对每个像素应用增强
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			r := float64(c.R)
			g := float64(c.G)
			b := float64(c.B)

			// Gamma 校正：gamma < 1 提亮暗部，gamma > 1 压暗暗部
			// 墨水屏暗部容易糊，用 0.8 略微提亮暗部
			const gamma = 0.8
			r = math.Pow(r/255.0, gamma) * 255.0
			g = math.Pow(g/255.0, gamma) * 255.0
			b = math.Pow(b/255.0, gamma) * 255.0

			// 对比度拉伸
			r = stretchValue(r, minGray, maxGray)
			g = stretchValue(g, minGray, maxGray)
			b = stretchValue(b, minGray, maxGray)

			// 饱和度增强：帮助墨水屏区分红/黄与黑白
			grayMid := r*0.299 + g*0.587 + b*0.114
			const satFactor = 1.5
			r = grayMid + (r-grayMid)*satFactor
			g = grayMid + (g-grayMid)*satFactor
			b = grayMid + (b-grayMid)*satFactor

			dst.SetRGBA(x, y, color.RGBA{
				uint8(clamp(int(r))),
				uint8(clamp(int(g))),
				uint8(clamp(int(b))),
				c.A,
			})
		}
	}

	return dst
}

// 对比度拉伸
func stretchValue(v, min, max float64) float64 {
	if max <= min {
		return v
	}
	return (v - min) * 255 / (max - min)
}

// 将图像编码为 BMP 格式（4色索引，墨水屏专用）
func encodeToBMP(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// 4 种墨水屏支持的颜色（调色板）
	// 顺序：黑、白、红、黄（索引 0,1,2,3）
	palette := []color.RGBA{
		{0, 0, 0, 255},       // 0 - 黑色
		{255, 255, 255, 255}, // 1 - 白色
		{255, 0, 0, 255},     // 2 - 红色
		{255, 255, 0, 255},   // 3 - 黄色
	}

	// 查找颜色索引的辅助函数（使用感知距离匹配，避免精确匹配失败）
	getColorIndex := func(c color.RGBA) uint8 {
		bestIdx := uint8(0)
		minDist := int(0x7FFFFFFF)
		for i, p := range palette {
			d := colorDistance(c, p)
			if d < minDist {
				minDist = d
				bestIdx = uint8(i)
			}
		}
		return bestIdx
	}

	// BMP 文件头
	var buf bytes.Buffer

	// 文件类型 "BM"
	buf.Write([]byte{0x42, 0x4D})

	// 文件大小（暂设为 0，稍后计算）
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00})

	// 保留字段
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00})

	// 位图数据偏移量
	// 文件头(14) + 信息头(40) + 调色板(4*4) = 14+40+16 = 70
	buf.Write([]byte{0x46, 0x00, 0x00, 0x00})

	// 位图信息头（40字节）
	buf.Write([]byte{0x28, 0x00, 0x00, 0x00}) // 信息头大小(40)

	// 宽度（小端字节序）
	buf.WriteByte(uint8(width & 0xFF))
	buf.WriteByte(uint8((width >> 8) & 0xFF))
	buf.WriteByte(uint8((width >> 16) & 0xFF))
	buf.WriteByte(uint8((width >> 24) & 0xFF))

	// 高度（小端字节序）
	buf.WriteByte(uint8(height & 0xFF))
	buf.WriteByte(uint8((height >> 8) & 0xFF))
	buf.WriteByte(uint8((height >> 16) & 0xFF))
	buf.WriteByte(uint8((height >> 24) & 0xFF))

	// 颜色平面数（1）
	buf.Write([]byte{0x01, 0x00})

	// 每像素位数（4 位，因为 4 种颜色需要 2 位，但 4 位更常用）
	// 注意：墨水屏可能需要 2 位，但 4 位兼容性更好
	buf.Write([]byte{0x04, 0x00})

	// 压缩方式（无压缩）
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00})

	// 位图数据大小（先设为 0）
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00})

	// 水平分辨率（像素/米，设为 2835）
	buf.Write([]byte{0x13, 0x0B, 0x00, 0x00})

	// 垂直分辨率（像素/米，设为 2835）
	buf.Write([]byte{0x13, 0x0B, 0x00, 0x00})

	// 调色板颜色数（4 种）
	buf.Write([]byte{0x04, 0x00, 0x00, 0x00})

	// 重要颜色数（4 种）
	buf.Write([]byte{0x04, 0x00, 0x00, 0x00})

	// 调色板（每个颜色 BGR + 0）
	for _, c := range palette {
		buf.WriteByte(c.B)  // B
		buf.WriteByte(c.G)  // G
		buf.WriteByte(c.R)  // R
		buf.WriteByte(0x00) // 保留字节
	}

	// 像素数据（4位像素，行对齐到 4 字节边界）
	// 每行的字节数：((width * 4 + 31) / 32) * 4
	rowSize := ((width*4 + 31) / 32) * 4 // 对齐到 4 字节边界

	// 从下到上写入像素数据（BMP 格式特点）
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		var rowBuf bytes.Buffer
		var currentByte uint8 = 0
		pixelCount := 0

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			idx := getColorIndex(c)

			// 每 2 个像素组成一个字节（4位/像素）
			// 注意：墨水屏有些是低4位在前，有些是高4位在前
			// 这里假设：第一个像素在高4位，第二个在低4位
			if pixelCount%2 == 0 {
				currentByte = (idx & 0x0F) << 4
			} else {
				currentByte |= (idx & 0x0F)
				rowBuf.WriteByte(currentByte)
			}
			pixelCount++
		}

		// 如果像素数为奇数，写入最后一个字节
		if pixelCount%2 != 0 {
			rowBuf.WriteByte(currentByte)
		}

		// 补齐到 4 字节对齐
		padding := make([]byte, rowSize-rowBuf.Len())
		rowBuf.Write(padding)

		buf.Write(rowBuf.Bytes())
	}

	// 更新文件大小
	fileSize := buf.Len()
	buf.Bytes()[2] = uint8(fileSize & 0xFF)
	buf.Bytes()[3] = uint8((fileSize >> 8) & 0xFF)
	buf.Bytes()[4] = uint8((fileSize >> 16) & 0xFF)
	buf.Bytes()[5] = uint8((fileSize >> 24) & 0xFF)

	return buf.Bytes(), nil
}

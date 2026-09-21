package tts

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SynthesizeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSynthesizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SynthesizeLogic {
	return &SynthesizeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SynthesizeLogic) Synthesize(req *types.TtsReq) (resp *types.TtsResp, err error) {
	textLen := utf8.RuneCountInString(req.Text)
	if textLen == 0 || textLen > 60 {
		return nil, fmt.Errorf("text 长度需 1~60 字 (当前 %d)", textLen)
	}

	per := clampInt(req.Per, 0, 5000)
	spd := clampInt(req.Spd, 0, 15)
	pit := clampInt(req.Pit, 0, 15)
	vol := clampInt(req.Vol, 0, 15)

	audio, err := l.svcCtx.BaiduTTS.Synthesize(l.ctx, req.Text, per, spd, pit, vol)
	if err != nil {
		return nil, fmt.Errorf("百度 TTS 失败: %w", err)
	}
	if len(audio) == 0 {
		return nil, errors.New("百度 TTS 返回空音频")
	}

	return &types.TtsResp{Audio: audio}, nil
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
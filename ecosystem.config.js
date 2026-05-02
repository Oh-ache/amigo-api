module.exports = {
  apps : [
    {
      name: "baseCodeRpc",
      script: "go",
      args: "run baseCode.go",
      cwd: "./app/baseCode/rpc/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/baseCodeRpc.log",
      time: true
    },
    {
      name: "baseCodeApi",
      script: "go",
      args: "run baseCode.go",
      cwd: "./app/baseCode/api/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/baseCodeApi.log",
      time: true
    },
    {
      name: "sdkRpc",
      script: "go",
      args: "run sdk.go",
      cwd: "./app/sdk/rpc/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/sdkRpc.log",
      time: true
    },
    {
      name: "sdkApi",
      script: "go",
      args: "run sdk.go",
      cwd: "./app/sdk/api/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/sdkApi.log",
      time: true
    },
    {
      name: "userRpc",
      script: "go",
      args: "run user.go",
      cwd: "./app/user/rpc/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/userRpc.log",
      time: true
    },
    {
      name: "userApi",
      script: "go",
      args: "run user.go",
      cwd: "./app/user/api/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/userApi.log",
      time: true
    },
    {
      name: "deviceRpc",
      script: "go",
      args: "run device.go",
      cwd: "./app/device/rpc/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/deviceRpc.log",
      time: true
    },
    {
      name: "deviceApi",
      script: "go",
      args: "run device.go",
      cwd: "./app/device/api/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/deviceApi.log",
      time: true
    },
    {
      name: "aiRpc",
      script: "go",
      args: "run ai.go",
      cwd: "./app/ai/rpc/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/aiRpc.log",
      time: true
    },
    {
      name: "aiApi",
      script: "go",
      args: "run ai.go",
      cwd: "./app/ai/api/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/aiApi.log",
      time: true
    },
    // {
      // name: "queue",
      // script: "go",
      // args: "run queue.go",
      // cwd: "./app/job/queue/",
      // instances: 1,
      // exec_mode: "fork",
      // watch: false,
      // log_file: "~/.logs/queue.log",
      // time: true
    // },
    {
      name: "mqueue",
      script: "go",
      args: "run mqueue.go",
      cwd: "./app/job/mqueue/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/mqueue.log",
      time: true
    },
    {
      name: "gateway",
      script: "go",
      args: "run gateway.go",
      cwd: "./app/gateway/",
      instances: 1,
      exec_mode: "fork",
      watch: false,
      log_file: "~/.logs/api.log",
      time: true
    },
  ],

  deploy : {
    production : {
      user : 'SSH_USERNAME',
      host : 'SSH_HOSTMACHINE',
      ref  : 'origin/master',
      repo : 'GIT_REPOSITORY',
      path : 'DESTINATION_PATH',
      'pre-deploy-local': '',
      'post-deploy' : 'npm install && pm2 reload ecosystem.config.js --env production',
      'pre-setup': ''
    }
  }
};

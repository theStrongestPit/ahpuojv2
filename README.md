### AHPUOJv2
#### 项目结构
```
├── config // 项目配置文件
├── controller 
├── core   // hustoj判题机源码
├── docker // 用于部署开发环境
├── docker-deploy // 用于部署线上环境
├── middleware // 中间件
├── migrate // 迁移脚本
├── model // 数据模型
├── request // 表单验证模型
├── router // 后台路由
├── service // mysql 和 redis 服务
├── static // 前端静态资源
├── tools // 开发时使用的工具函数
├── utils // 工具函数
├── vendor // go依赖
├── web-admin // 后台前端
├── web-common // 前端公共
└── web-user // 前台前端
├── webpack.config.js
├── webpack.prod.config.js
├── main.go
├── index.html // 前台入口文件
├── admin_index.html // 后台管理系统入口文件
├── package.json
├── README.md
```
#### 部署开发环境
* go 的依赖采用 govendor 管理，保存在 vendor 目录下
* 在项目目录下
```
npm install // 安装npm依赖
npm run dev // 启动开发模式（热加载）
```
* cd 进入 docker 目录，使用命令启动容器
```
docker-compose up -d
```
* 进入容器，开发环境做了项目目录的映射，项目目录映射到了 容器内 /home/judge/go/src/ahpuoj目录
```
docker-compose exec ahpuoj bash
```
* 在容器内进入 /home/judge/go/src/ahpuoj/docker-deploy/ahpuojv2-deploy
/install 目录，运行脚本 install.sh，设置容器内环境
* 在容器 /home/judge/go/src/ahpuoj 目录下运行命令
```
gowatch
```
* 开发环境部署完成，建议使用 vscode-insider 的 remote-docker 扩展，在容器内开发，便于 go 程序的调试

#### 服务器部署
* 将 docker-deploy 目录上传到服务器
* 同开发环境部署类似，启动容器，服务器环境做了容器 web目录 和 data目录(保存题目数据) 的映射
* 确保各项配置都已经正确设置，使用 go build 生成 go 的可执行程序，运行 npm run build 生成前端文件
* 参考 remote_deploy_mycloud.sh 将程序上传到服务器上映射的 容器 web 目录，即部署完成
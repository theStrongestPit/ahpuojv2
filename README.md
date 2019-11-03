### AHPUOJv2

#### 项目结构

```
├── config // 项目配置文件
├── controller
├── docker
├───core 判题机源码
├───dev 开发环境容器
└───prod 生产环境容器
├── middleware // 中间件
├── model // 数据模型
├── request // 表单验证模型
├── router // 后台路由
├── service // mysql 和 redis 服务
├── tools // 开发时使用的工具函数
├── utils // 工具函数
├── vendor // go依赖
├── web-admin // 后台管理系统前端
├── web-common // 前端公共
└── web-user // oj前台前端
├── webpack.base.js // 公用 webpack 配置
├── webpack.config.js // 开发环境 webpack 配置
├── webpack.prod.config.js // 生产环境构建 webpack 配置
├── main.go // go 程序入口
```

#### 准备工作

因为 OJ 系统的判题机需要运行在一个比较复杂的环境中，所以部署工作较为复杂

- 将 docker/dev 和 docker/prod 目录中的 .env.example 文件复制一份，重命名为 .env，并将空缺的参数填写上去
- 将 config 目录下的 config.ini.example 复制一份，重命名为 config.ini，并填上参数，参数设置应该与 docker/dev 和 /docker/prod 中的 .env 文件中的参数一致
- docker/prod/ahpuojv2-deploy/assets/install 中的 config.ini.example 文件也做同样的处理

#### 部署开发环境

- go 的依赖采用 govendor 管理，保存在 vendor 目录下
- 在项目目录下运行以下命令，启动前端的开发环境

```
npm install // 安装npm依赖
npm run dev // 启动开发模式
```

- cd 进入 docker/dev 目录，使用命令启动容器

```
docker-compose up -d
```

- 进入容器，开发环境做了项目目录的映射，项目目录映射到了 容器内 /home/judge/go/src/ahpuoj 目录

```
docker-compose exec ahpuoj bash
```

- 在容器内进入 /home/judge/go/src/ahpuoj/docker/dev 目录，运行脚本 install.sh
- 在容器 /home/judge/go/src/ahpuoj 目录下运行命令，可以在开发时对 go 程序进行热编译

```
gowatch
```

- 开发环境部署完成，建议使用 vscode-insider 的 remote-docker 扩展，在容器内开发，便于 go 程序的调试

#### 服务器部署

- 将 docker 目录整体上传到服务器
- cd 进入 docker/prod 目录，使用命令启动容器，服务器环境做了容器 web 目录（存放前端资源） 、 data 目录（保存题目数据）和 core 目录（判题机源码）的 的映射
- 确保各项配置都已经正确设置，在本机项目文件夹下使用 go build 生成 go 的可执行程序，运行 npm run build 生成前端文件
- 参考 remote_deploy_mycloud.sh 将程序上传到服务器上映射的 容器 web 目录
- 在服务器进入部署环境容器，使用 supervisorctl 运行 ahpuoj

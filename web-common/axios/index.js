import axios from 'axios';
import Cookies from 'js-cookie';
// import { Message } from 'element-ui'

let baseURL = '/api';

// 添加请求拦截器
axios.interceptors.request.use(
  config => {
    // 如果本地有token
    if (Cookies.get('access-token')) {
      config.headers = {
        Authorization: Cookies.get('access-token')
      };
    }
    return config;
  },
  err => {
    return Promise.reject(err);
  }
);

// // 添加响应拦截器
// axios.interceptors.response.use((res) => {
//     // 对响应数据做点什么
//     return res;
// }, (err) => {
//     return Promise.reject(err);
// });

// 封装数据返回失败提示函数
function errorState(response) {
  // 如果http状态码正常，则直接返回数据
  if (
    response &&
    (response.status === 200 ||
      response.status === 304 ||
      response.status === 400)
  ) {
    return response;
  } else {
    // Message({
    //     message: '服务器内部错误',
    //     type: 'error'
    // });
    return response;
  }
}

// 封装数据返回成功提示函数
// function successState(res) {
//     // 统一判断后端返回的错误码(错误码与后台协商而定)
//     if (res.data.code === '000000') {
//         console.log('success')
//         return res
//     }
// }

// 封装axios
function request(method, url, payload) {
  let httpDefault = {
    method: method,
    baseURL: baseURL,
    url: url,
    // `params` 是即将与请求一起发送的 URL 参数
    // `data` 是作为请求主体被发送的数据
    params: method === 'GET' || method === 'DELETE' ? payload : null,
    data: method === 'POST' || method === 'PUT' ? payload : null,
    timeout: 10000
  };

  return new Promise(async (resolve, reject) => {
    try {
      let res = await axios(httpDefault);
      // successState(res)
      resolve(res);
    } catch (err) {
      errorState(err);
      reject(err);
    }
  });
}

export default request;

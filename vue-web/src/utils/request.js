import axios from 'axios'; // 引入axios
import { Message } from 'element-ui';
import { store } from '@/store'
import context from '@/main.js'
const service = axios.create({
    baseURL: process.env.VUE_APP_BASE_API,
    timeout: 99999
})
let acitveAxios = 0
let timer
const showLoading = () => {
    acitveAxios++
    if (timer) {
        clearTimeout(timer)
    }
    timer = setTimeout(() => {
        if (acitveAxios > 0) {
            context.$bus.emit("showLoading")
        }
    }, 400);
}

const closeLoading = () => {
        acitveAxios--
        if (acitveAxios <= 0) {
            clearTimeout(timer)
            context.$bus.emit("closeLoading")
        }
    }
    //http request 拦截器
service.interceptors.request.use(
    config => {
        if (!config.donNotShowLoading) {
            showLoading()
        }
        const token = store.getters['user/token']
        const user = store.getters['user/userInfo']
        config.data = JSON.stringify(config.data);
        config.headers = {
            'Content-Type': 'application/json',
            'x-token': token,
            'x-user-id': user.id,
            'x-user-name': user.userName,
        }
        return config;
    },
    error => {
        closeLoading()
        Message({
            showClose: true,
            message: error,
            type: 'error'
        })
        return error;
    }
);


//http response 拦截器
service.interceptors.response.use(
    response => {
        closeLoading()
        console.log(response)
        if (response.headers["new-token"]) {
            store.commit('user/setToken', response.headers["new-token"])
        }
        if (response.data.code == 0) {
            console.log(response.data)
            return response.data
        } else {
            Message({
                showClose: true,
                message: response.data.msg || decodeURI(response.headers.msg),
                type: 'error',
            })
            if (response.data.data && response.data.data.reload) {
                console.log("resp:",response.data.data)
                store.commit('user/LoginOut')
            }
            return response.data.msg ? response.data : response
        }
    },
    error => {
        closeLoading()
        Message({
            showClose: true,
            message: error,
            type: 'error'
        })
        return error
    }
)

export default service
import service from '@/utils/request'

export const login = (data) => {
    return service({
        url: "/user/login",
        method: 'post',
        data: data
    })
}
export const getUserList = (data) => {
    return service({
        url: "/user/getUserList",
        method: 'post',
        data: data
    })
}
export const setUserRole = (data) => {
    return service({
        url: "/user/setUserRole",
        method: 'post',
        data: data
    })
}
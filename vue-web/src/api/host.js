import service from '@/utils/request'

export const getHostList = (data) => {
    return service({
        url: "/api/getHostList",
        method: 'post',
        data
    })
}

export const createHost = (data) => {
    return service({
        url: "/api/createHost",
        method: 'post',
        data
    })
}

export const getHostById = (data) => {
    return service({
        url: "/api/getHostById",
        method: 'post',
        data
    })
}

export const checkHost = (data) => {
    return service({
        url: "/api/checkHost",
        method: 'post',
        data
    })
}

export const reloadHost = (data) => {
    return service({
        url: "/api/reloadHost",
        method: 'post',
        data
    })
}
export const deleteHost = (data) => {
    console.log('request.data:',data)
    return service({
        url: "/api/deleteHost",
        method: 'post',
        data
    })
}
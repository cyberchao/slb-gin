import service from '@/utils/request'

export const getUpstreamList = (data) => {
    return service({
        url: "/api/getUpstreamList",
        method: 'post',
        data
    })
}

export const createUpstream = (data) => {
    return service({
        url: "/api/createUpstream",
        method: 'post',
        data
    })
}

export const updateUpstream = (data) => {
    return service({
        url: "/api/updateUpstream",
        method: 'post',
        data
    })
}

export const deleteUpstream = (data) => {
    console.log('request.data:',data)
    return service({
        url: "/api/deleteUpstream",
        method: 'post',
        data
    })
}
import axios from 'axios'

const isDev = import.meta.env.DEV 
const baseUrl = isDev ? 'http://localhost:3000/items' : '/api/items'

const getAll = async () => {
    const request = await axios.get(baseUrl)
    return request.data
}

export default { getAll }
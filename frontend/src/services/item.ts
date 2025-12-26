import axios from 'axios'
const baseUrl = 'http://localhost:3000/items'

const getAll = async () => {
    const request = await axios.get(baseUrl)
    return request.data
}

export default { getAll }
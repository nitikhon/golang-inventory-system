package util

func GetOffset(page, limit int) int {
    if page < 1 { page = 1 }
    if limit < 1 { limit = 12 }
    return (page - 1) * limit
}
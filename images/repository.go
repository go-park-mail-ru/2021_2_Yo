package images

type Repository interface {
	StoreImg(userId string, fileName string)
}
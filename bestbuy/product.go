package bestbuy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/couchbase/gocb.v1"
)

type Product struct {
	Sku               string   `json:"sku"`
	ID                string   `json:"id"`
	Type              string   `json:"type"`
	Name              string   `json:"name"`
	RegPrice          float64  `json:"reg_price"`
	SalePrice         float64  `json:"sale_price"`
	Discount          float64  `json:"discount"`
	OnSale            string   `json:"on_sale"`
	ShortDescription  string   `json:"short_description"`
	Class             string   `json:"class"`
	BbItemID          string   `json:"bb_item_id"`
	ModelNumber       string   `json:"model_number"`
	Manufacturer      string   `json:"manufacturer"`
	Image             string   `json:"image"`
	MedImage          string   `json:"med_image"`
	ThumbImage        string   `json:"thumb_image"`
	LargeImage        string   `json:"large_image"`
	LongDescription   string   `json:"long_description"`
	Keywords          string   `json:"keywords"`
	CatDescendentPath string   `json:"cat_descendent_path"`
	CatIds            []string `json:"cat_ids"`
}

func CbConnect(address string, username string, password string) (*gocb.Cluster, error) {
	cluster, err := gocb.Connect(address)
	if err != nil {
		return nil, err
	}
	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func CbOpenBucket(name string, cluster *gocb.Cluster) (*gocb.Bucket, error) {
	bucket, err := cluster.OpenBucket(name, "")
	return bucket, err
}

func LoadProductsFromFile(pathToJson string) error {

	// open the bucket
	cluster, err := gocb.Connect("couchbase://localhost")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "evan",
		Password: "password",
	})
	if err != nil {
		return err
	}
	bucket, err := cluster.OpenBucket("bb-catalog", "")
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(pathToJson)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(file)
	for {
		s, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		product := Product{}
		err = json.Unmarshal([]byte(s), &product)
		if err != nil {
			return err
		}
		err = SaveProduct(product, bucket)
		if err != nil {
			return err
		}
	}
	err = bucket.Close()
	return err
}

func SaveProduct(product Product, bucket *gocb.Bucket) error {
	_, err := bucket.Upsert(product.ID, product, 0)
	return err
}

func GetProduct(ID string, bucket *gocb.Bucket) (*Product, error) {
	var product Product
	_, err := bucket.Get(ID, &product)
	return &product, err
}

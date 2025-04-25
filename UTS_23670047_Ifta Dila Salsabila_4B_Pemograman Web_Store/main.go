package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Kategori struct {
	ID           uint   `gorm:"primaryKey"`
	NamaKategori string `gorm:"column:nama_kategori"`
}

func (Kategori) TableName() string {
	return "kategori"
}

type Produk struct {
	ID         uint    `gorm:"primaryKey"`
	Nama       string  `gorm:"column:nama"`
	Stok       int     `gorm:"column:stok"`
	Harga      float64 `gorm:"column:harga"`
	KategoriID uint    `gorm:"column:kategori_id"`
	Kategori   Kategori
}

func (Produk) TableName() string {
	return "produk"
}

func connectDB() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/store?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// HOME / READ
func homeHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectDB()

	var produks []Produk
	var kategoris []Kategori

	// Preload kategori untuk produk
	db.Preload("Kategori").Find(&produks)
	db.Find(&kategoris)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Produks":   produks,
		"Kategoris": kategoris,
	})
}

// CREATE & UPDATE
func produkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, _ := strconv.Atoi(r.FormValue("id")) // Cek apakah ada id untuk edit
		nama := r.FormValue("nama")
		stok, _ := strconv.Atoi(r.FormValue("stok"))
		harga, _ := strconv.ParseFloat(r.FormValue("harga"), 64)
		kategoriID, _ := strconv.Atoi(r.FormValue("kategori_id"))

		db, _ := connectDB()

		// Jika id == 0 berarti tambah produk, jika id > 0 berarti edit produk
		if id == 0 {
			// Menambah produk baru
			db.Create(&Produk{
				Nama:       nama,
				Stok:       stok,
				Harga:      harga,
				KategoriID: uint(kategoriID),
			})
		} else {
			// Mengupdate produk yang ada
			db.Model(&Produk{}).Where("id = ?", id).Updates(Produk{
				Nama:       nama,
				Stok:       stok,
				Harga:      harga,
				KategoriID: uint(kategoriID),
			})
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE
func hapusProdukHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	db, _ := connectDB()
	db.Delete(&Produk{}, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ================== MAIN ==================
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/produk", produkHandler) // Menangani baik tambah maupun edit
	http.HandleFunc("/hapus", hapusProdukHandler)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	log.Println("Server running at http://localhost:8083")
	http.ListenAndServe(":8083", nil)
}

package main

import (
	"api-go/cmd/auth"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gonum.org/v1/gonum/mat"
)

type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type QRResponse struct {
	Q [][]float64 `json:"Q"`
	R [][]float64 `json:"R"`
}

func QRFactorization(matrix [][]float64) ([][]float64, [][]float64, error) {
	rows := len(matrix)
	cols := len(matrix[0])

	data := make([]float64, 0, rows*cols)
	for i := 0; i < rows; i++ {
		data = append(data, matrix[i]...)
	}

	A := mat.NewDense(rows, cols, data)

	var qr mat.QR
	qr.Factorize(A)

	Q := mat.NewDense(rows, rows, nil)
	qr.QTo(Q)

	R := mat.NewDense(rows, cols, nil)
	qr.RTo(R)

	qData := make([][]float64, rows)
	rData := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		qData[i] = make([]float64, rows)
		rData[i] = make([]float64, cols)
		for j := 0; j < rows; j++ {
			qData[i][j] = Q.At(i, j)
		}
		for j := 0; j < cols; j++ {
			rData[i][j] = R.At(i, j)
		}
	}

	return qData, rData, nil
}

func ProtectedQRHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req MatrixRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al decodificar la matriz", http.StatusBadRequest)
		return
	}

	if len(req.Matrix) == 0 || len(req.Matrix[0]) == 0 {
		http.Error(w, "Matriz vacía no es válida", http.StatusBadRequest)
		return
	}

	Q, R, err := QRFactorization(req.Matrix)
	if err != nil {
		http.Error(w, "Error al realizar la factorización QR", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(QRResponse{
		Q: Q,
		R: R,
	})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.Handle("/qr", auth.JWTMiddleware(http.HandlerFunc(ProtectedQRHandler))).Methods("POST")

	log.Println("Servidor iniciado en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

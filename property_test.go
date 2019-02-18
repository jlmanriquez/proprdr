package main

import "testing"

const (
	fileName = "./config.properties"
)

func TestNew(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("Error creando nuevo PropertyFile")
	}

	if propertyFile.Size() == 0 {
		t.Error("Error en lectura de PropertyFile. Se esperaban al menos un elemento")
	}

}

func TestGet(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("Error creando nuevo PropertyFile")
	}

	property := propertyFile.Get("app.name")
	if property != "Lector de Propiedades" {
		t.Errorf("Error en lectura de PropertyFile. Se esperaba 'Lector de Propiedades' y se encontro '%s'", property)
	}
}

func TestGetAsInt(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("Error creando nuevo PropertyFile")
	}

	intValue, err := propertyFile.GetAsInt("app.values.maxConnections")
	if err != nil {
		t.Errorf("Se esperaba obtener un entero pero se obtuvo error %s", err.Error())
	}

	if intValue < 20 {
		t.Errorf("Se esperaba un entero mayor o igual 20 y se obtuvo %d", intValue)
	}
}

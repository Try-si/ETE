package ETEHelper

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// LoadTMX charge et parse un fichier .tmx
func LoadTMX(path string) (*TMXMap, error) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture fichier %s: %v", path, err)
	}

	var mapData TMXMap
	if err := xml.Unmarshal(fileData, &mapData); err != nil {
		return nil, fmt.Errorf("erreur parse XML: %v", err)
	}

	// Post-traitement optionnel : parser les données CSV immédiatement pour gagner du temps à l'exécution
	if err := processLayers(&mapData); err != nil {
		return nil, err
	}

	return &mapData, nil
}

// LoadTSX charge et parse un fichier .tsx
func LoadTSX(path string) (*TSXTileset, error) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tileset TSXTileset
	if err := xml.Unmarshal(fileData, &tileset); err != nil {
		return nil, err
	}

	return &tileset, nil
}

// processLayers itère sur les calques et chunks pour décoder le CSV
func processLayers(m *TMXMap) error {
	for i := range m.Layers {
		if err := decodeLayerData(&m.Layers[i]); err != nil {
			return err
		}
		// Si le calque a des chunks (carte infinie)
		for j := range m.Layers[i].Chunks {
			if err := decodeChunkData(&m.Layers[i].Chunks[j]); err != nil {
				return err
			}
		}
	}

	// Gérer aussi les groupes si nécessaire
	for i := range m.Groups {
		for j := range m.Groups[i].Layers {
			if err := decodeLayerData(&m.Groups[i].Layers[j]); err != nil {
				return err
			}
			for k := range m.Groups[i].Layers[j].Chunks {
				if err := decodeChunkData(&m.Groups[i].Layers[j].Chunks[k]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// decodeLayerData parse le contenu CSV d'un calque standard
func decodeLayerData(layer *TMXLayer) error {
	if layer.Data.Encoding != "csv" {
		return fmt.Errorf("encodage %s non supporté pour le calque %s (utilisez csv)", layer.Data.Encoding, layer.Name)
	}

	// Le contenu brut contient des sauts de ligne et des virgules
	// On veut obtenir une slice []uint32 (ou uint64) de tous les GIDs
	gids, err := parseCSV(layer.Data.RawContent)
	if err != nil {
		return err
	}

	// Ici, vous pouvez stocker 'gids' dans un champ personnalisé de votre structure
	// ou les traiter directement. Pour l'exemple, on suppose que vous les gardez en mémoire.
	// Note: TMXLayerData ne contient pas de champ 'Tiles' par défaut, vous devrez peut-être
	// ajouter un champ `Tiles []uint32` dans votre struct TMXLayerData ou TMXLayer.
	// Pour cet exemple, je vais supposer que vous avez ajouté ce champ :
	// layer.Data.Tiles = gids
	_ = gids // À assigner à votre structure finale

	return nil
}

// decodeChunkData parse le CSV d'un chunk (carte infinie)
func decodeChunkData(chunk *TMXChunk) error {
	if chunk.Data.Encoding != "csv" {
		return fmt.Errorf("encodage non supporté dans un chunk")
	}

	gids, err := parseCSV(chunk.Data.RawContent)
	if err != nil {
		return err
	}

	// C'est ICI que l'information de position (chunk.X, chunk.Y) est cruciale.
	// Comme chunk.X et chunk.Y sont des INT, ils peuvent être négatifs.
	// Votre moteur devra calculer la position globale :
	// PosX = (chunk.X * chunk.Width) + offset_colonne
	// PosY = (chunk.Y * chunk.Height) + offset_ligne
	_ = gids
	return nil
}

// parseCSV transforme une chaîne "0,0,1,2,\n3,4..." en []uint32
func parseCSV(raw string) ([]uint32, error) {
	// Nettoyer les espaces, sauts de ligne
	clean := strings.ReplaceAll(raw, "\n", ",")
	clean = strings.ReplaceAll(clean, " ", "")
	parts := strings.Split(clean, ",")

	var tiles []uint32
	for _, p := range parts {
		if p == "" {
			continue
		}
		val, err := strconv.ParseUint(p, 10, 64) // Parse en 64 pour sécurité
		if err != nil {
			return nil, fmt.Errorf("erreur parse tuile '%s': %v", p, err)
		}
		// On cast en uint32 car un GID tient sur 32 bits (même avec les flags)
		tiles = append(tiles, uint32(val))
	}
	return tiles, nil
}

package store

import (
	"mywallet/utils"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func AddToken(id uint64, tokenDetails *utils.TokenDetails) error {
	// for converting int64 to time.Time  - Unix
	accessTime := time.Unix(tokenDetails.AccessExpiry, 0)
	refreshTime := time.Unix(tokenDetails.RefreshExpiry, 0)

	accessErr := Client.Set(tokenDetails.AccessID, strconv.Itoa(int(id)), accessTime.Sub(accessTime)).Err()
	if accessErr != nil {
		log.Error().Err(accessErr).Any("accessid", tokenDetails.AccessID).Any("access time", accessTime).
			Msg("error  in access for adding token")
		return accessErr
	}

	refreshErr := Client.Set(tokenDetails.RefreshID, strconv.Itoa(int(id)), accessTime.Sub(refreshTime)).Err()
	if refreshErr != nil {
		log.Error().Err(refreshErr).Any("refreshid", tokenDetails.RefreshID).Any("refresh time", refreshTime).
			Msg("error in refresh for adding token")
		return refreshErr
	}
	return nil
}

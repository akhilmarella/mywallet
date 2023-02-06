package store

import (
	"mywallet/api"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func AddToken(id uint64, tokenDetails *api.TokenDetails) error {
	// for converting int64 to time.Time  - Unix
	accessTime := time.Unix(tokenDetails.AccessExpiry, 0)
	refreshTime := time.Unix(tokenDetails.RefreshExpiry, 0)

	accessErr := Client.Set(tokenDetails.AccessID, strconv.Itoa(int(id)), accessTime.Sub(accessTime)).Err()
	if accessErr != nil {
		log.Error().Err(accessErr).Any("access_id", tokenDetails.AccessID).Any("access_time", accessTime).
			Msg("error  in access for adding token")
		return accessErr
	}

	refreshErr := Client.Set(tokenDetails.RefreshID, strconv.Itoa(int(id)), accessTime.Sub(refreshTime)).Err()
	if refreshErr != nil {
		log.Error().Err(refreshErr).Any("refresh_id", tokenDetails.RefreshID).Any("refresh_time", refreshTime).
			Msg("error in refresh for adding token")
		return refreshErr
	}
	return nil
}

func DeleteRefreshID(refreshID string) (int64, error) {
	deleted, err := Client.Del(refreshID).Result()
	if err != nil {
		log.Error().Err(err).Any("id", deleted).Msg("error in id")
		return 0, err
	}
	return deleted, nil
}

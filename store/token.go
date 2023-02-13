package store

import (
	"mywallet/api"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func AddToken(id int64, tokenDetails *api.TokenDetails) error {
	// for converting int64 to time.Time  - Unix
	accessTime := time.Unix(tokenDetails.AccessExpiry, 0)
	refreshTime := time.Unix(tokenDetails.RefreshExpiry, 0)

	accessErr := Client.Set(tokenDetails.AccessID, strconv.Itoa(int(id)), accessTime.Sub(accessTime)).Err()
	if accessErr != nil {
		log.Error().Err(accessErr).Any("access_id", tokenDetails.AccessID).Any("access_time", accessTime).
			Any("action:", "store_token.go_AddToken").
			Msg("error  in access for adding token")
		return accessErr
	}

	refreshErr := Client.Set(tokenDetails.RefreshID, strconv.Itoa(int(id)), accessTime.Sub(refreshTime)).Err()
	if refreshErr != nil {
		log.Error().Err(refreshErr).Any("refresh_id", tokenDetails.RefreshID).Any("refresh_time", refreshTime).
			Any("action:", "store_token.go_AddToken").
			Msg("error in refresh for adding token")
		return refreshErr
	}
	return nil
}

func DeleteRefreshID(refreshID string) (int64, error) {
	deleted, err := Client.Del(refreshID).Result()
	if err != nil {
		log.Error().Err(err).Any("id", deleted).
			Any("action:", "store_token.go_DeleteRefreshID").
			Msg("error in deleting refresh id")
		return 0, err
	}
	return deleted, nil
}

func DeleteTokens(accessID string, id int64) (int64, int64, error) {
	refreshID := accessID + "_" + strconv.Itoa(int(id))

	deleteAccessID, err := Client.Del(accessID).Result()
	if err != nil || deleteAccessID == 0 {
		log.Error().Err(err).Any("id", deleteAccessID).
			Any("action:", "store_token.go_DeleteTokens").
			Msg("error in deleting access id")
		return deleteAccessID, 0, err
	}

	deleteRefreshID, err := Client.Del(refreshID).Result()
	if err != nil || deleteRefreshID == 0 {
		log.Error().Err(err).Any("id", deleteRefreshID).
			Any("action:", "store_token.go_DeleteTokens").
			Msg("error in delete refresh id")
		return 0, deleteRefreshID, err
	}

	return deleteAccessID, deleteRefreshID, nil
}

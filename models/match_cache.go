package models

import (
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type MatchCache struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 承兑人账号
	Acceptor string `json:"acceptor" gorm:"column:acceptor"`
	// 卡id
	CardId int64 `json:"card_id" gorm:"column:card_id"`
	// 卡类型
	CardType string `json:"card_type" gorm:"column:card_type"`
	// 最大承兑金额
	MaxMatchedAmount int64 `json:"max_matched_amount" gorm:"column:max_matched_amount"`
	// 最小承兑金额
	MinMatchedAmount int64 `json:"min_matched_amount" gorm:"column:min_matched_amount"`
	// 上一次的匹配时间
	LastMatchTime util.JSONTime `json:"last_match_time" gorm:"column:last_match_time"`
	// 资金版本号
	FundVersion int64 `json:"fund_version" gorm:"column:fund_version"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置MatchCache的表名为`match_cache`
func (MatchCache) TableName() string {
	return "match_cache"
}

func NewMatchCacheModel(session *Session) *MatchCache {
	return &MatchCache{Session: session}
}

// ExistMatchCacheByID checks if an matchCache exists based on ID
func (a *MatchCache) ExistMatchCacheByID(id int64) (bool, error) {
	matchCache := new(MatchCache)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(matchCache).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if matchCache.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetMatchCacheTotal gets the total number of matchCaches based on the constraints
func (a *MatchCache) GetMatchCacheTotal(maps interface{}) (int, error) {
	var count int
	err := a.Session.db.Model(&MatchCache{}).Where(maps).Count(&count).Error
	return count, err
}

// GetMatchCaches gets a list of matchCaches based on paging constraints
func (a *MatchCache) MatchOneCardCache(maps interface{}) (*MatchCache, error) {
	matchCache := new(MatchCache)
	err := a.Session.db.Where(maps).First(matchCache).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return matchCache, nil
}

// GetMatchCache Get a single matchCache based on ID
func (a *MatchCache) GetMatchCache(id int64) (*MatchCache, error) {
	matchCache := new(MatchCache)
	err := a.Session.db.Where("id = ?", id).First(matchCache).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return matchCache, nil
}

// EditMatchCache modify a single matchCache
func (a *MatchCache) EditMatchCache(id int64, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&MatchCache{}).Where("id = ?", id).Updates(data).Error
}

// AddMatchCache add a single matchCache
func (a *MatchCache) AddMatchCache(matchCache *MatchCache) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(matchCache).Error
}

// DeleteMatchCache delete a single matchCache
func (a *MatchCache) DeleteMatchCache(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(MatchCache{}).Error
}

//更新缓存卡最大匹配金额
func (a *MatchCache) UpdateMatchedCardMaxMatchedAmount(id int64, maxMatchedAmount int64) error {
	tx := GetSessionTx(a.Session)

	data := map[string]interface{}{
		"max_matched_amount": maxMatchedAmount,
		"last_match_time":    util.JSONTimeNow(),
	}
	return tx.Model(&MatchCache{}).Where("id = ?", id).Updates(data).Error
}

func (a *MatchCache) GetMatchedCard(agency string, amount int64, cardType string) (*MatchCache, error) {
	matchCache := new(MatchCache)

	err := a.Session.db.Where("agency = ? and  card_type = ?", agency, cardType).Where("max_matched_amount >= ? "+
		"and min_matched_amount <= ?", amount, amount).Select("card_id").Order("last_match_time").First(matchCache).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return matchCache, nil
}

func GetMatchedCard(agency string, amount int64, cardType string) ([]int64, error) {
	var matchCaches []*MatchCache
	var matchedCards []int64

	err := db.Where("agency = ? and  card_type = ?", agency, cardType).Where("max_matched_amount >= ? "+
		"and min_matched_amount <= ?", amount, amount).Select("card_id").Order("last_match_time").Find(&matchCaches).Error

	if err != nil {
		return matchedCards, err
	}
	for _, matchCache := range matchCaches {
		matchedCards = append(matchedCards, matchCache.CardId)
	}
	return matchedCards, nil
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	ds "github.com/star-find-cloud/star-mall/internal/deepseek"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/repo"
	"strconv"
)

type DeepseekService interface {
	// SuggestProductByUserTags 根据用户标签，获取与用户标签匹配的热门商品
	SuggestProductByUserTags(ctx context.Context, userid int64, model domain.DeepseekModel) (<-chan string, error)

	// SearchProductSuggest 根据用户搜索商品, 进行商品推荐
	SearchProductSuggest(ctx context.Context, searchMsg string, model domain.DeepseekModel) (<-chan string, error)
}

type DeepseekServiceImpl struct {
	dsAdapter   *ds.DeepSeekAdapter
	productRepo repo.ProductRepo
	userRepo    repo.UserRepo
}

func NewDeepseekService(dsAdapter *ds.DeepSeekAdapter, userRepo *repo.UserRepoImpl, productRepo *repo.ProductRepoImpl) *DeepseekServiceImpl {
	return &DeepseekServiceImpl{dsAdapter: dsAdapter, userRepo: userRepo, productRepo: productRepo}
}

// SuggestProductByUserTags 根据用户标签，获取与用户标签匹配的热门商品
func (s *DeepseekServiceImpl) SuggestProductByUserTags(ctx context.Context, userid int64, model domain.DeepseekModel) (<-chan string, error) {
	// 获取用户近期感兴趣的标签
	userTags, err := s.userRepo.GetUserTags(ctx, userid)
	if err != nil {
		applog.AppLogger.Errorf("can't get user's tags: %v", err)
		return nil, fmt.Errorf("can't get user's tags: %w", err)
	}
	// 过滤掉空标签
	if len(userTags) == 0 {
		return nil, fmt.Errorf("user's tags is empty")
	}

	var tagsMap map[string]string
	err = json.Unmarshal(userTags, &tagsMap)
	if err != nil {
		applog.AppLogger.Errorf("can't unmarshal user's tags: %v", err)
		return nil, fmt.Errorf("can't unmarshal user's tags: %w", err)
	}

	errSlice := make([]string, 0)
	tagsStrSlice := make([]string, 0)

	var tags = make([]int64, 0)
	// 将用户标签转换为产品类型
	for _, tagStr := range tagsMap {
		//fmt.Println("tagStr", tagStr)
		tag, err := strconv.ParseInt(tagStr, 10, 64)
		//fmt.Println("tag", tag)
		if err != nil {
			errSlice = append(errSlice, fmt.Sprintf("标签格式错误: %s", tagStr))
		}
		value, exist := _const.ProductTypeMap[tag]
		if !exist {
			errSlice = append(errSlice, fmt.Sprintf("未存在的标签: %d", tag))
		}
		tagsStrSlice = append(tagsStrSlice, value)
		tags = append(tags, tag)
	}
	if len(errSlice) > 0 {
		return nil, fmt.Errorf("不能找到与标签相匹配的商品: %v", errSlice)
	}

	//fmt.Println(tags)
	// 获取与用户标签匹配的热门商品
	products, err := s.productRepo.GetByCateIDsAndHot(ctx, tags)
	//fmt.Println(products)

	// 将 商品 和 用户标签 转换为 DeepSeek 可以识别的格式
	var msg = s.dsAdapter.GenerateUserMsg(products, tagsStrSlice)
	//fmt.Println("msg:", msg)

	// 生成消息通道, 供 DeepSeek 回复消息
	var messageChan <-chan string
	// 根据用户选择的模型，生成回复消息
	if model == domain.ModelChat {
		messageChan, err = s.dsAdapter.GenerateReplyStreamChat(ctx, msg, 2, string(domain.ModelChat))
		if err != nil {
			applog.AppLogger.Errorf("can't generate reply stream chat: %v", err)
			return nil, fmt.Errorf("can't generate reply stream chat: %w", err)
		}
	} else if model == domain.ModelResearcher {
		messageChan, err = s.dsAdapter.GenerateReplyStreamReasoner(ctx, msg, 2, string(domain.ModelResearcher))
		if err != nil {
			applog.AppLogger.Errorf("can't generate reply stream reasoner: %v", err)
			return nil, fmt.Errorf("can't generate reply stream reasoner: %w", err)
		}
	}

	return messageChan, err
}

// SearchProductSuggest 根据用户搜索商品, 进行商品推荐
func (s *DeepseekServiceImpl) SearchProductSuggest(ctx context.Context, searchMsg string, model domain.DeepseekModel) (<-chan string, error) {
	products, err := s.productRepo.SearchByMsg(ctx, searchMsg)
	if err != nil {
		applog.AppLogger.Errorf("can't search product by msg: %v", err)
		return nil, fmt.Errorf("can't search product by msg: %w", err)
	}
	var msg = s.dsAdapter.GenerateUserMsg(products, searchMsg)
	var messageChan <-chan string
	if model == domain.ModelChat {
		messageChan, err = s.dsAdapter.GenerateReplyStreamChat(ctx, msg, 2, string(domain.ModelChat))
	} else if model == domain.ModelResearcher {
		messageChan, err = s.dsAdapter.GenerateReplyStreamReasoner(ctx, msg, 2, string(domain.ModelResearcher))
	}

	return messageChan, err
}

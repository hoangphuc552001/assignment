package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/store/memory"
	"net/http"
	"sfvn-test/common/model"
	"sfvn-test/common/util"
	"sfvn-test/service"
	"time"
)

type Histories struct {
	histories service.IHistories
}

func NewAPIHistories(r *gin.Engine, histories service.IHistories) {
	handler := &Histories{
		histories: histories,
	}
	limiterRate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  1,
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, limiterRate)
	Group := r.Group("v1/get-histories")
	{
		Group.GET("", util.RateLimiter(rateLimiter), handler.GetHistories)
	}
}

func (h *Histories) GetHistories(context *gin.Context) {
	startDate := context.Query("start_date")
	endDate := context.Query("end_date")
	period := context.Query("period")
	symbol := context.Query("symbol")
	var startDateParse, endDateParse time.Time
	startDateParse, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "start_date is invalid"})
		return
	}
	endDateParse, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "end_date is invalid"})
		return
	}
	if startDateParse.After(endDateParse) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be after end_date"})
		return
	}
	if symbol == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	var periodEnum model.Period
	if period == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "period is required"})
		return
	} else {
		periodEnum = model.Period(period)
		if periodEnum != model.Period1D && periodEnum != model.Period1H && periodEnum != model.Period30M {
			context.JSON(http.StatusBadRequest, gin.H{"error": "period is invalid"})
			return
		}
	}
	historiesParams := model.HistoriesRequest{
		StartDate: startDateParse,
		EndDate:   endDateParse,
		Period:    periodEnum,
		Symbol:    symbol,
	}
	histories, err := h.histories.GetHistories(historiesParams)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, histories)
}

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"sfvn-test/common/model"
	"testing"
)

type mockHistories struct{}

func (m *mockHistories) GetHistories(params model.HistoriesRequest) ([]model.HistoriesResponse, error) {
	return []model.HistoriesResponse{}, nil
}

func TestNewAPIHistories(t *testing.T) {
	r := gin.New()
	histories := &mockHistories{}
	NewAPIHistories(r, histories)
	req, err := http.NewRequest("GET", "/v1/get-histories?start_date=2022-01-01&end_date=2022-01-10&period=1D&symbol=BTC", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetHistories(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?start_date=2022-01-01&end_date=2022-01-10&period=1D&symbol=BTC", nil)
	h.GetHistories(context)
	assert.Equal(t, http.StatusOK, context.Writer.Status())
}

func TestGetHistoriesWithInvalidStartDate(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?start_date=2022-01-01&period=1D&symbol=BTC", nil)
	h.GetHistories(context)
	assert.NotEqual(t, http.StatusOK, context.Writer.Status())
}

func TestGetHistoriesWithInvalidEndDate(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?end_date=2022-01-10&period=1D", nil)
	h.GetHistories(context)
	assert.NotEqual(t, http.StatusOK, context.Writer.Status())
}

func TestGetHistoriesWithInvalidStartDateAfterEndDate(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?start_date=2022-03-03&end_date=2022-01-10&period=1D", nil)
	h.GetHistories(context)
	assert.NotEqual(t, http.StatusOK, context.Writer.Status())
}

func TestGetHistoriesWithEmptySymbol(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?start_date=2022-03-03&end_date=2022-04-10&period=1D", nil)
	h.GetHistories(context)
	assert.NotEqual(t, http.StatusOK, context.Writer.Status())
}

func TestGetHistoriesWithEmptyPeriod(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?start_date=2022-03-03&end_date=2022-04-10&symbol=bitcoin", nil)
	h.GetHistories(context)
	assert.NotEqual(t, http.StatusOK, context.Writer.Status())
}

func TestGetHistoriesWithInvalidPeriod(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?start_date=2022-03-03&end_date=2022-04-10&period=1A&symbol=bitcoin", nil)
	h.GetHistories(context)
	assert.NotEqual(t, http.StatusOK, context.Writer.Status())
}

func TestGetHistoriesWrongLogic(t *testing.T) {
	h := &Histories{
		histories: &mockHistories{},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest("GET", "/v1/get-histories?start_date=2022-03-01&end_date=2022-01-01&period=1D&symbol=bitcoin", nil)
	h.GetHistories(context)
	assert.NotEqual(t, http.StatusOK, context.Writer.Status())
}

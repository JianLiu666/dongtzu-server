package server

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/repository/githubSDK"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"gitlab.geax.io/demeter/gologger/logger"
)

func getProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		ctx := context.Background()
		providerProfile, err := arangodb.GetProviderProfileByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(model.GetProviderInfoRes{
			StatusCode: "200",
			Data:       providerProfile,
		})
	}
}

func registerProvider() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. parsing post body
		var body model.RegisterProviderReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No validate registeration input",
			})
		}

		// 2. register provider tx
		ctx := context.Background()
		err := arangodb.CreateProviderProfile(ctx, body)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(model.RegisterOrUpdateProviderRes{
			StatusCode: "200001",
		})
	}
}

// Todo
// - 1. 之後再實作手機修改簡訊驗證邏輯
// - 2. 之後再實作gmail修改email + calendar同步邏輯
// 目前只實做完全信任前端覆蓋資料的邏輯
func updateProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		var body model.UpdateProviderInfoReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No validate edit provider input",
			})
		}

		// 2. update provider profile
		ctx := context.Background()
		err := arangodb.UpdateProviderByLineUserID(ctx, lineUserID, body)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		// 3. if params status is 2 && update success -> github 串接
		if body.Status == constant.Provider_Status_Auditing {
			profile, _ := arangodb.GetProviderProfileByLineUserID(ctx, lineUserID)
			err = githubSDK.CreateIssueForProvider(*profile)
			if err != nil {
				logger.Errorf("[githubSDK] create issue failure : %v", err)
			}
		}

		return c.Status(fasthttp.StatusOK).JSON(model.RegisterOrUpdateProviderRes{
			StatusCode: "200001",
		})
	}
}

// 拿Provider班表
func getProviderEventSchedule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		ctx := context.Background()
		schedulesWithCourse, err := arangodb.GetScheduleByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		appointments, err := arangodb.GetAppointmentsByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		res := formatEventSchedule(schedulesWithCourse, appointments)

		return c.Status(fasthttp.StatusOK).JSON(res)
	}
}

// 拿Provider收入簡介
func getProviderIncomeSummary() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		ctx := context.Background()

		payments, err := arangodb.GetPaymentsByLineUserID(ctx, lineUserID, 2)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		receipts, err := arangodb.GetMonthReceiptByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		res := formatIncomeSummaryFromPaymentsAndReceipts(payments, receipts)

		return c.Status(fasthttp.StatusOK).JSON(res)
	}
}

// 歷史紀錄/分頁
func getMonthReceipts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		// validate inputs
		limit := c.Query("limit", "20")
		perPage, err := strconv.Atoi(limit)
		if err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404002",
				StatusCode: "404002",
				Message:    "Not validate input",
			})
		}
		page := c.Query("page", "1")
		currentPage, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404002",
				StatusCode: "404002",
				Message:    "Not validate input",
			})
		}

		queryStart := c.Query("start", "")
		start, _ := strconv.Atoi(queryStart)

		ctx := context.Background()
		receipts, err := arangodb.GetMonthReceiptList(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		res := formatMonthReceiptList(receipts, perPage, currentPage, start)

		return c.Status(fasthttp.StatusOK).JSON(res)
	}
}

func getProviderServiceProducts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		ctx := context.Background()
		svcProducts, err := arangodb.GetServiceProductsByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(model.ProviderServiceProductListRes{
			StatusCode: "200",
			Data:       svcProducts,
		})
	}
}

func createOrUpdateProviderServiceProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		var body model.CreateOrUpdateServiceProductsReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No validate edit provider input",
			})
		}

		// 2. update provider profile
		ctx := context.Background()
		err := arangodb.CreateOrUpdateServiceProduct(ctx, lineUserID, body)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(model.CreateUpdateDeleteRes{
			StatusCode: "200001",
		})
	}
}

// 拿未來兩個月的班表
func getProviderSchedule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		// validate inputs
		limit := c.Query("limit", "20")
		perPage, err := strconv.Atoi(limit)
		if err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404002",
				StatusCode: "404002",
				Message:    "Not validate input",
			})
		}
		page := c.Query("page", "1")
		currentPage, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404002",
				StatusCode: "404002",
				Message:    "Not validate input",
			})
		}

		queryStart := c.Query("start", "")
		start, _ := strconv.Atoi(queryStart)

		ctx := context.Background()
		schedules, err := arangodb.GetScheduleList(ctx, lineUserID, int64(start))
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		res := formatScheduleList(schedules, perPage, currentPage)

		return c.Status(fasthttp.StatusOK).JSON(res)
	}
}

// 創建單堂預約
func createServiceSchedule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		var body model.CreateServiceScheduleReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No validate edit provider input",
			})
		}

		ctx := context.Background()
		err := arangodb.CreateProviderSchedule(ctx, lineUserID, body)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(model.CreateUpdateDeleteRes{
			StatusCode: "200001",
		})
	}
}

// 創建課程規則
func createScheduleRule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		// Todo
		// 1. body parser
		// 2. create single schedule
		//    - check if the schedule rule is conflict
		//    - create the schedule rule

		return nil
	}
}

/**
 * format method
 */
func formatIncomeSummaryFromPaymentsAndReceipts(payments []model.Payment,
	receipts []model.MonthReceipt) model.ProviderIncomeSummaryRes {
	var res model.ProviderIncomeSummaryRes

	var sumOfPayments int
	for _, payment := range payments {
		sumOfPayments += int(payment.NetAmount)
	}

	sort.Slice(receipts, func(i, j int) bool {
		return receipts[j].ClearingStartedAt > receipts[i].ClearingStartedAt
	})

	res.StatusCode = string(fasthttp.StatusOK)
	res.Data = model.IncomeSummary{
		SumOfPayments:      sumOfPayments,
		LastMonthIncome:    int(receipts[1].NetIncome),
		CurrentMonthIncome: int(receipts[0].NetIncome),
	}
	return res
}

func formatEventSchedule(schedules []model.ScheduleCourse,
	appointments []model.Appointment) model.ProviderEventScheduleRes {
	var res model.ProviderEventScheduleRes
	var events []model.EventSchedule

	for _, schedule := range schedules {
		events = append(events, model.EventSchedule{
			Start:        schedule.CourseStartAt,
			End:          schedule.CourseEndAt,
			ScheduleType: 0,
			Title:        schedule.Title,
			Content:      schedule.Content,
			Count:        schedule.Count,
		})
	}

	for _, appt := range appointments {
		events = append(events, model.EventSchedule{
			Start:        appt.CourseStartAt,
			End:          appt.CourseEndAt,
			ScheduleType: 1,
			Title:        "",
			Content:      "",
			Count:        0,
		})
	}

	res.StatusCode = string(fasthttp.StatusOK)
	res.Data = events
	return res
}

func formatMonthReceiptList(receipts []model.MonthReceipt,
	perPage, curPage, start int) model.ProviderMonthReceiptList {
	var res model.ProviderMonthReceiptList
	var list []model.MonthReceipt

	if start > 0 {
		for _, receipt := range receipts {
			if receipt.ClearingStartedAt >= int64(start) {
				list = append(list, receipt)
			}
		}
	}

	sort.Slice(receipts, func(i, j int) bool {
		return receipts[j].ClearingStartedAt > receipts[i].ClearingStartedAt
	})
	totalCounts, totalPages, items := paginate(list, perPage, curPage)

	res.StatusCode = string(fasthttp.StatusOK)
	res.Data = items
	res.Meta = model.Pagination{
		TotalCounts: totalCounts,
		TotalPages:  totalPages,
		CurrentPage: curPage,
		PerPage:     perPage,
	}
	return res
}

func paginate(inputs []model.MonthReceipt, perPage int, currentPage int) (int, int, []model.MonthReceipt) {
	totalCounts := len(inputs)
	if totalCounts == 0 {
		return 0, 0, nil
	}

	lastPageNumber := perPage
	totalPages := totalCounts / perPage
	lastPageCount := totalCounts % perPage
	if lastPageCount != 0 {
		totalPages++
		if currentPage == totalPages {
			lastPageNumber = lastPageCount
		}
	}

	start := (currentPage - 1) * perPage
	end := start + lastPageNumber

	return totalCounts, totalPages, inputs[start:end]
}

func formatScheduleList(schedules []model.Schedule, perPage, curPage int) model.ProviderScheduleListRes {
	var res model.ProviderScheduleListRes
	var list []model.Schedule

	sort.Slice(schedules, func(i, j int) bool {
		return schedules[j].CourseStartAt > schedules[i].CourseStartAt
	})
	totalCounts, totalPages, items := paginateSchedule(list, perPage, curPage)

	res.StatusCode = string(fasthttp.StatusOK)
	res.Data = items
	res.Meta = model.Pagination{
		TotalCounts: totalCounts,
		TotalPages:  totalPages,
		CurrentPage: curPage,
		PerPage:     perPage,
	}
	return res
}

func paginateSchedule(inputs []model.Schedule, perPage int, currentPage int) (int, int, []model.Schedule) {
	totalCounts := len(inputs)
	if totalCounts == 0 {
		return 0, 0, nil
	}

	lastPageNumber := perPage
	totalPages := totalCounts / perPage
	lastPageCount := totalCounts % perPage
	if lastPageCount != 0 {
		totalPages++
		if currentPage == totalPages {
			lastPageNumber = lastPageCount
		}
	}

	start := (currentPage - 1) * perPage
	end := start + lastPageNumber

	return totalCounts, totalPages, inputs[start:end]
}

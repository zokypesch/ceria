// func getRatePlans(c client.Client, ctx context.Context, hotelIDs []string, inDate string, outDate string) chan rateResults {
// 	rateClient := rate.NewRateService("go.micro.srv.rate", c)
// 	ch := make(chan rateResults, 1)

// 	go func() {
// 		res, err := rateClient.GetRates(ctx, &rate.Request{
// 			HotelIds: hotelIDs,
// 			InDate:   inDate,
// 			OutDate:  outDate,
// 		})
// 		ch <- rateResults{res.RatePlans, err}
// 	}()

// 	return ch

// how to call
// rateCh := getRatePlans(s.Client, ctx, nearby.HotelIds, inDate, outDate)

// // wait on rates reply
// rateReply := <-rateCh
// if err := rateReply.err; err != nil {
// 	return merr.InternalServerError("api.hotel.rates", err.Error())
// }
// }

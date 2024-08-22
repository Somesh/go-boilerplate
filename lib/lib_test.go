package lib

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"gopkg.in/tokopedia/logging.v1"
)

func TestMain(m *testing.M) {
	flag.Parse()
	logging.SetDebug(testing.Verbose()) // turn debug on if verbose mode

	os.Exit(m.Run())
}

func TestTimeToUnixSeconds(t *testing.T) {
	format := "2006-01-02"
	date := "1983-10-15"

	unixSeconds, err := TimeToUnixSeconds(context.Background(), date, format, "")
	if err != nil {
		t.Error(err)
	}
	if unixSeconds != 435024000 {
		t.Error("Unable to convert date to unixSeconds ", unixSeconds)
	}

	unixSeconds, err = TimeToUnixSeconds(context.Background(), date, "", "")
	if err != nil {
		t.Error(err)
	}
	if unixSeconds != 435024000 {
		t.Error("Unable to convert date to unixSeconds ", unixSeconds)
	}

	unixSeconds, err = TimeToUnixSeconds(context.Background(), date, "", "Invalid")
	if err == nil {
		t.Errorf("Expcted error , got nil")
	}

	date = "90"
	unixSeconds, err = TimeToUnixSeconds(context.Background(), date, "", "Invalid")
	if err == nil {
		t.Errorf("Expcted error , got nil")
	}
}

func TestUnixSecondsToDate(t *testing.T) {
	format := "2006-01-02 15:04:05"
	date := int64(435024000)

	ctime := UnixSecondsToDate(context.Background(), date, format)
	if ctime != "1983-10-15 00:00:00" {
		t.Error("Unable to convert unix seconds to date time format")
	}
	ctime = UnixSecondsToDate(context.Background(), date, "")
	if ctime != "1983-10-15 00:00:00" {
		t.Error("Unable to convert unix seconds to date time format")
	}
}

func TestNormalizePrice(t *testing.T) {
	priceString := "123.456"
	priceInt64 := int64(35000)
	priceFloat := 123.456
	priceFloat2Place := 123.45

	normalizedString, errString := NormalizePrice(context.Background(), priceString)
	if normalizedString != int64(12345) {
		t.Error("Invalid value for string", normalizedString)
	}
	if errString != nil {
		t.Error(errString.Error())
	}
	normalizedInt64, errInt := NormalizePrice(context.Background(), priceInt64)
	if normalizedInt64 != int64(3500000) {
		t.Error("Invalid value for int64", normalizedInt64)
	}
	if errInt != nil {
		t.Error(errInt.Error())
	}
	normalizedFloat, errFloat := NormalizePrice(context.Background(), priceFloat)
	if normalizedFloat != int64(12345) {
		t.Error("Invalid value for float", normalizedFloat)
	}
	if errFloat != nil {
		t.Error(errFloat.Error())
	}
	normalizedFloat2place, errFloat2Place := NormalizePrice(context.Background(), priceFloat2Place)
	if normalizedFloat != int64(12345) {
		t.Error("Invalid value for float 2 place", normalizedFloat2place)
	}
	if errFloat2Place != nil {
		t.Error(errFloat2Place.Error())
	}
}

func TestRawstringToInterface(t *testing.T) {
	input := `{"name": "somesh"}`
	result, err := RawstringToInterface(context.Background(), input)
	if err != nil {
		t.Error("Invalid JSON ", input, " Err: ", err)
	}
	t.Log(result)

	input = `{"somesh"}`
	result, err = RawstringToInterface(context.Background(), input)
	if err == nil {
		t.Errorf("Expected error , got nil")
	}

	input = `"somesh"'}`
	result, err = RawstringToInterface(context.Background(), input)
	if err == nil {
		t.Errorf("Expected error , got nil")
	}
}

func TestFormatNumber(t *testing.T) {
	inputs := []interface{}{
		"123",
		"4205",
		"20",
		"8000",
		50000,
		"1000000",
		"30000000",
		800000000,
		-10000,
		4205.0,
	}

	outputs := []string{
		"123",
		"4.205",
		"20",
		"8.000",
		"50.000",
		"1.000.000",
		"30.000.000",
		"800.000.000",
		"-10.000",
		"4.205",
	}

	for idx, x := range inputs {
		val := FormatNumber(context.Background(), x)
		t.Log(inputs[idx], "equivalent", val)
		if val != outputs[idx] {
			t.Error("Failed for ", x)
		}
	}
}

func TestQueryParamToIntArr(t *testing.T) {
	input := "123,4,5,6,7,8,9,asd"
	output := []int64{123, 4, 5, 6, 7, 8, 9, 0}

	result := QueryParamToIntArr(input)

	for idx, x := range result {
		if output[idx] != x {
			t.Error("Failed for ", x)
		}
	}
}

func TestStringArrayToMap(t *testing.T) {
	input := []string{"somesh", "jitan", "ishaan", "shivani"}
	output := true
	result := StringArrayToMap(context.Background(), input)
	if result["jitan"] != output {
		t.Error("Test fail to convert string array to map ")
	}
}

func TestTimeByZone(t *testing.T) {
	expiry, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "2018-10-17 03:44:18 +0000 UTC")
	result := "17 Oct 2018 10:44"

	resp := TimeByZone(context.Background(), "", "02 Jan 2006 15:04", expiry)
	if resp != result {
		t.Error("Test fail to convert time by zone", resp)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	val := []int64{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 3}
	result := []int64{1, 2, 3, 4, 5}

	res := RemoveDuplicates(val)

	var totalRes, totalResp int64
	totalResp = 15
	if len(res) == len(result) {
		for idx := range res {
			totalRes += res[idx]
		}
		if totalRes != totalResp {
			t.Error("Test fail to remove duplicate integers", res)
		}
	} else {
		t.Error("Test fail to remove duplicate integers", res)
	}

}

func TestIntArrayToMap(t *testing.T) {
	ctx := context.Background()
	val := []int64{1, 2, 3, 4, 5}
	var total int64
	total = 15

	resp := IntArrayToMap(ctx, val)
	var respTotal int64
	for key := range resp {
		respTotal += key
	}

	if total != respTotal {
		t.Error("Test failed in converting the int array to map", resp)
	}
}

func TestConvertDefaultTimeToTimeZone(t *testing.T) {
	ctx := context.Background()
	utc := time.Now().UTC()
	locJkt := time.FixedZone("Asia/Jakarta", 7*3600)
	locInd := time.FixedZone("Asia/Calcutta", 5.5*3600)

	utcTime, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", utc.String())
	currTime := utcTime.In(locInd)

	result := utcTime.In(locJkt)

	updatedTime, err := ConvertTimeToTimeZone(ctx, currTime, "Asia/Jakarta", 7*3600)
	if err != nil {
		t.Error("Test Failed to convert time to jkt time zone", err)
	}
	if updatedTime.String() != result.String() {
		t.Error("Wrong time convertion", updatedTime, result, utcTime, currTime)
	}
}

func TestGetDateRange(t *testing.T) {
	ctx := context.Background()
	s1 := int64(123)
	e1 := int64(456)
	result := GetDateRange(ctx, s1, e1)
	if result {
		t.Errorf("Expected False. Got True")
	}

	e1 = 456000000
	result = GetDateRange(ctx, s1, e1)
	if !result {
		t.Errorf("Expected True , Got False")
	}
}

func TestFormatProoviderID(t *testing.T) {
	ctx := context.Background()
	id := "1"
	categoryID := int64(456)
	providerID := int64(999)
	expectedResult := "1-456-999"

	result := FormatProviderId(ctx, id, categoryID, providerID)
	if result != expectedResult {
		t.Errorf("Expected %s , Got %s", expectedResult, result)
	}
}

func TestGetProviderCtxId(t *testing.T) {
	ctx := context.Background()
	localID := "123-456-789"
	expectedResult := "123"
	result := GetProviderCtxId(ctx, localID)
	if result != expectedResult {
		t.Errorf("Expected %s , Got %s", expectedResult, result)
	}

	localID = ""
	expectedResult = ""
	result = GetProviderCtxId(ctx, localID)
	if result != expectedResult {
		t.Errorf("Expected %s , Got %s", expectedResult, result)
	}
}

func TestMakeProductURL(t *testing.T) {
	ctx := context.Background()
	name := "test"
	id := int64(321)
	expectedResult := "test-321"
	result := MakeProductURL(ctx, name, id)
	if result != expectedResult {
		t.Errorf("Expected %s , Got %s", expectedResult, result)
	}
}

func TestUrlizeString(t *testing.T) {
	ctx := context.Background()
	input := "test url"
	expectedResult := "test-url"
	result := UrlizeString(ctx, input)
	if result != expectedResult {
		t.Errorf("Expected %s , Got %s", expectedResult, result)
	}
}

func TestGetColumnName(t *testing.T) {
	ctx := context.Background()
	type InputStruct struct {
		name string `db:"name"`
	}
	input := InputStruct{name: "Test"}
	result := GetColumnName(ctx, &input, "name")
	if result != "name" {
		t.Errorf("Got %s", result)
	}
}

func TestToString(t *testing.T) {
	ctx := context.Background()
	i1 := int64(12)
	i2 := 99
	i3 := "q1w2"
	i4 := 45.2
	expected1 := "12"
	expected2 := "99"
	expected3 := "q1w2"
	expected4 := "45"

	r1 := ToString(ctx, i1)
	if r1 != expected1 {
		t.Errorf("Expected %s , Got %s", expected1, r1)
	}
	r2 := ToString(ctx, i2)
	if r2 != expected2 {
		t.Errorf("Expected %s , Got %s", expected2, r2)
	}
	r3 := ToString(ctx, i3)
	if r3 != expected3 {
		t.Errorf("Expected %s , Got %s", expected3, r3)
	}
	r4 := ToString(ctx, i4)
	if r4 != expected4 {
		t.Errorf("Expected %s , Got %s", expected4, r4)
	}
}

func TestDenormalizePrice(t *testing.T) {
	ctx := context.Background()
	input := int64(1000)
	expected := int64(10)
	result := DeNormalizePrice(ctx, input)
	if result != expected {
		t.Errorf("Expected %d, Got %d", expected, result)
	}
}

func TestToInt64(t *testing.T) {
	ctx := context.Background()
	i1 := int64(12)
	i2 := "99"
	i3 := "q1w2"
	i4 := 45.2
	expected1 := int64(12)
	expected2 := int64(99)
	expected3 := int64(0)
	expected4 := int64(45)

	r1 := ToInt64(ctx, i1)
	if r1 != expected1 {
		t.Errorf("Expected %d , Got %d", expected1, r1)
	}
	r2 := ToInt64(ctx, i2)
	if r2 != expected2 {
		t.Errorf("Expected %d , Got %d", expected2, r2)
	}
	r3 := ToInt64(ctx, i3)
	if r3 != expected3 {
		t.Errorf("Expected %d , Got %d", expected3, r3)
	}
	r4 := ToInt64(ctx, i4)
	if r4 != expected4 {
		t.Errorf("Expected %d , Got %d", expected4, r4)
	}

}

func TestIntArrayToString(t *testing.T) {
	ctx := context.Background()
	input := []int64{1, 2, 3}
	delimier := ","
	expected := "1,2,3"
	result := IntArrayToString(ctx, input, delimier)
	if result != expected {
		t.Errorf("Expected %s , Got %s", expected, result)
	}
}

func TestValidatePromocode(t *testing.T) {
	ctx := context.Background()
	promo1 := "AB"
	promo2 := "TEST_PROMO"

	result := ValidatePromocode(ctx, promo1)
	if result {
		t.Errorf("Expected False , Got True")
	}
	result = ValidatePromocode(ctx, promo2)
	if !result {
		t.Errorf("Expected True , Got False")
	}
}

func TestGeneratePDF(t *testing.T) {
	ctx := context.Background()
	input := []byte{10, 20}
	_, _ = GeneratePDF(ctx, input)
}

func TestGetOutBoundIP(t *testing.T) {
	_, _ = GetOutboundIP()
}

func TestValidateStatus(t *testing.T) {
	ctx := context.Background()
	s1 := []int64{1, 2, 3}
	s2 := []int64{4, 5, -5, 6}

	result := ValidateStatus(ctx, s1)
	if !result {
		t.Errorf("Expected True , Got False")
	}
	result = ValidateStatus(ctx, s2)
	if result {
		t.Errorf("Expected False, Got True")
	}
}

func TestTruncateString(t *testing.T) {
	ctx := context.Background()
	input := "1234567890"
	maxLength := 5
	expected := "12345"
	result := TruncateString(ctx, input, maxLength)
	if result != expected {
		t.Errorf("Expected %s , Got %s", expected, result)
	}

}

func TestTimeBeginningOfWeek(t *testing.T) {
	month := time.Month(11)
	now := time.Date(2009, month, 10, 23, 0, 0, 0, time.UTC)
	t1 := TimeBeginningOfWeek(now, false)
	t2 := TimeBeginningOfWeek(now, true)
	expected1 := time.Date(2009, month, 9, 0, 0, 0, 0, time.UTC)
	expected2 := time.Date(2009, month, 8, 0, 0, 0, 0, time.UTC)
	if t1 != expected1 {
		t.Errorf("Expected %+v, Got %+v", expected1, t1)
	}
	if t2 != expected2 {
		t.Errorf("Expected %+v, Got %+v", expected2, t2)
	}
}

func TestTimeEndOfWeek(t *testing.T) {
	month := time.Month(11)
	now := time.Date(2009, month, 10, 23, 0, 0, 0, time.UTC)
	t1 := TimeEndOfWeek(now, false)
	t2 := TimeEndOfWeek(now, true)
	expected1 := time.Date(2009, month, 15, 23, 59, 59, 999999999, time.UTC)
	expected2 := time.Date(2009, month, 14, 23, 59, 59, 999999999, time.UTC)
	if t1 != expected1 {
		t.Errorf("Expected %+v, Got %+v", expected1, t1)
	}
	if t2 != expected2 {
		t.Errorf("Expected %+v, Got %+v", expected2, t2)
	}
}

func TestTimeBeginningOfMonth(t *testing.T) {
	month := time.Month(11)
	now := time.Date(2009, month, 10, 23, 0, 0, 0, time.UTC)
	t1 := TimeBeginningOfMonth(now)
	expected1 := time.Date(2009, month, 1, 0, 0, 0, 0, time.UTC)
	if t1 != expected1 {
		t.Errorf("Expected %+v, Got %+v", expected1, t1)
	}
	t1 = TimeEndOfMonth(now)
	expected1 = time.Date(2009, month, 30, 0, 0, 0, 0, time.UTC)
	if t1 != expected1 {
		t.Errorf("Expected %+v, Got %+v", expected1, t1)
	}
}

func TestZoneTimeNow(t *testing.T) {
	ctx := context.Background()
	loc := "UTC"
	now := time.Now()
	expected := now.Format("2006-01-02")
	result := ZoneTimeNow(ctx, loc, "")
	if result != expected {
		t.Errorf("Expected %s , Got %s", expected, result)
	}
}

func TestRawstringToSructure(t *testing.T) {
	ctx := context.Background()
	raw := `{"a":"1"`
	type dest struct {
		A string `json:"a"`
	}
	expected := dest{A: "1"}
	var d dest
	err := RawstringToStructure(ctx, raw, nil)
	if err == nil {
		t.Errorf("Expected error , got nil")
	}

	raw = `{"a":"1"}`
	err = RawstringToStructure(ctx, raw, &d)
	if err != nil {
		t.Error(err)
	}
	if d != expected {
		t.Errorf("Expected %+v , Got %+v", expected, d)
	}
}

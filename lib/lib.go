package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"

	"gopkg.in/tokopedia/logging.v1"

	"github.com/Somesh/go-boilerplate/common/constant"
)

var trimOutside = regexp.MustCompile(`^[\*\\\"\s\p{Zs}]+|[\*\\\"\s\p{Zs}]+$`)
var trimInside = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

func TimeToUnixSeconds(ctx context.Context, date, format, location string) (int64, error) {
	if format == "" {
		format = constant.DateFormat
	}

	if location == "" {
		location = constant.LOC_DEFAULT
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		return -1, err
	}

	t, err := time.ParseInLocation(format, date, loc)
	if err != nil {
		return -1, err
	}

	return t.UTC().UnixNano() / int64(time.Second), err
}

func GetDateRange(ctx context.Context, startSeconds, endSeconds int64) bool {
	timeDiff := endSeconds - startSeconds

	if timeDiff >= constant.Seconds24HRS {
		return true
	} else {
		return false
	}
}

func UnixSecondsToDate(ctx context.Context, seconds int64, format string) string {
	if format == "" {
		format = constant.DateTimeFormat
	}

	ctime := time.Unix(seconds, 0).UTC()
	return ctime.Format(format)
}

func FormatProviderId(ctx context.Context, Id string, categoryId, providerId int64) string {
	return fmt.Sprintf("%s-%d-%d", Id, categoryId, providerId)
}

func GetProviderCtxId(ctx context.Context, localID string) string {
	/* LocalId is the one which we store at our end by doing FormatProviderId
	This function returns the actual provider context id, it can be provider product_id,
	schedule_id section_id etc.
	*/
	return strings.Split(localID, "-")[0]
}

func MakeProductURL(ctx context.Context, name string, id int64) string {
	s := UrlizeString(ctx, name)
	return fmt.Sprintf("%s-%d", strings.ToLower(s), id)
}

func UrlizeString(ctx context.Context, input string) string {
	return strings.Replace(GetRunedString(ctx, input), " ", "-", -1)
}

func GetRunedString(ctx context.Context, input string) string {
	return stripCtlAndExtFromUnicode(ctx, input)
}

func GetColumnName(ctx context.Context, instance interface{}, structField string) string {
	f, ok := reflect.TypeOf(instance).Elem().FieldByName(structField)
	if !ok {
		log.Printf("[lib][GetColumnName] Error. StructField does not exists :%s", structField)
	}
	return string(f.Tag.Get("db"))
}

func GetFieldValue(ctx context.Context, instance interface{}, structField string) interface{} {
	r := reflect.ValueOf(instance)
	return reflect.Indirect(r).FieldByName(structField)
}

func ToString(ctx context.Context, input interface{}) string {
	var output string
	switch t := input.(type) {
	case int, int8, int64:
		output = fmt.Sprintf("%d", t)
	case string:
		output = fmt.Sprintf("%s", t)
	case float64:
		output = strconv.Itoa(int(t))
	default:
		resultJSON, err := json.Marshal(input)
		if err == nil {
			output = string(resultJSON)
		} else {
			log.Printf("expected int, string or float64 but got %s", reflect.TypeOf(input))
		}
	}
	return output
}

func NormalizePrice(ctx context.Context, input interface{}) (int64, error) {
	/*
		This function takes input as string, float, int, int64
		and returns corresponding value multiplied by 1000. This is done avoid storing
	*/

	var output int64
	var err error

	switch t := input.(type) {
	case string:
		in, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			fIn, err := strconv.ParseFloat(t, 64)
			if err != nil {
				break
			}
			output = FloatToInt64(ctx, fIn*100)
			break
		}
		output = in * 100
	case int, int64:
		output = t.(int64) * 100
	case float64:
		output = FloatToInt64(ctx, t*100)
	default:
		log.Printf("[NormalizePrice]expected int, string or float64 but got %s", reflect.TypeOf(input))
	}

	return output, err
}

func DeNormalizePrice(ctx context.Context, price int64) int64 {
	price = price / 100
	return (price * 100) / 100
}

// To denormalize to floating point value
func DeNormalizeToFloat(ctx context.Context, price int64) float64 {
	return float64(price) / 100
}

func FloatToInt64(ctx context.Context, in float64) int64 {
	return int64(in)
}

func ToInt64(ctx context.Context, input interface{}) int64 {
	var output int64
	switch t := input.(type) {
	case int, int8, int64:
		output = t.(int64)
	case string:
		output, _ = strconv.ParseInt(t, 10, 64)
	case float64:
		output = FloatToInt64(ctx, t)
	case reflect.Value:
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			output = int64(t.Int())
		case reflect.String:
			output, _ = strconv.ParseInt(t.String(), 10, 64)
		}
	default:
		log.Printf("[ToInt64]expected int, string or float64 but got %s, Input %+v", reflect.TypeOf(input), input)
	}
	return output
}

func IntArrayToString(ctx context.Context, ar []int64, delim string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ar)), delim), "[]")
}

func TrimString(ctx context.Context, source string) string {
	final := trimOutside.ReplaceAllString(source, "")
	final = trimInside.ReplaceAllString(final, " ")
	return final
}

func ValidatePromocode(ctx context.Context, promocode string) bool {
	match, _ := regexp.MatchString("[A-z0-9_]+$", promocode)
	if len(promocode) < 3 || len(promocode) > 30 || !match {
		return false
	}
	return true
}

// TimeBeginningOfWeek return the begin of the week of t
// bSundayFirst means that many country use the monday as the first day of week
func TimeBeginningOfWeek(t time.Time, bSundayFirst bool) time.Time {

	weekday := int(t.Weekday())
	if !bSundayFirst {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}

	d := time.Duration(-weekday) * 24 * time.Hour
	t = t.Add(d)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).UTC()
}

// TimeEndOfWeek return the end of the week of t
// bSundayFirst means that many country use the monday as the first day of week
func TimeEndOfWeek(t time.Time, bSundayFirst bool) time.Time {
	return TimeBeginningOfWeek(t, bSundayFirst).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// TimeBeginningOfMonth return the begin of the month of t
func TimeBeginningOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).UTC()
}

// TimeEndOfMonth return the end of the month of t
func TimeEndOfMonth(t time.Time) time.Time {
	return TimeBeginningOfMonth(t).AddDate(0, 1, -1)
}

// set offset in no of months
func DateByOffset(t time.Time, offset int) time.Time {
	return TimeEndOfMonth(t).AddDate(0, offset, 0)
}

// returs now for a given timezone

//format := "2006-01-02T15:04:05Z"

func ZoneTimeNow(ctx context.Context, location, format string) string {

	if format == "" {
		format = constant.DateFormat
	}

	loc, _ := time.LoadLocation(location)
	now := time.Now().In(loc).Format(format)

	return now
}

func RawstringToInterface(ctx context.Context, raw string) (map[string]interface{}, error) {
	var result map[string]interface{}

	rawIn := json.RawMessage(raw)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		logging.Debug.Println("[rawstringToInterface] Unable to Marshal input. Err: ", err)
		return nil, err
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		logging.Debug.Println("[rawstringToInterface] Unable to UnMarshal input. Err: ", err)
		return nil, err
	}

	return result, nil
}

func FormatNumber(ctx context.Context, input interface{}) string {

	v := reflect.ValueOf(input)
	var x string

	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		x = fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		x = fmt.Sprintf("%d", int(v.Float()))
	case reflect.String:
		x = v.String()
	}

	lastIndex := len(x) - 1

	var buffer []byte
	var strBuffer bytes.Buffer

	j := 0
	for i := lastIndex; i >= 0; i-- {
		j++
		buffer = append(buffer, x[i])

		if j == 3 && i > 0 {
			buffer = append(buffer, '.')
			j = 0
		}
	}

	for i := len(buffer) - 1; i >= 0; i-- {
		strBuffer.WriteByte(buffer[i])
	}
	result := strBuffer.String()
	return result
}

func GeneratePDF(ctx context.Context, inputBytes []byte) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		logging.Debug.Printf("[lib][GeneratePDF] Error in generating new Pdf object. Error: %+v", err)
		return nil, err
	}

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(inputBytes))

	chErr := make(chan error, 1)

	go func() {
		pdfg.AddPage(page)

		if err = pdfg.Create(); err != nil {
			chErr <- err
			return
		}
		chErr <- nil
	}()

	select {
	case err := <-chErr:
		if err != nil {
			logging.Debug.Println("[lib][GeneratePDF] Error in creating new Pdf")
			return nil, err
		}

	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return pdfg.Bytes(), nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}

func QueryParamToIntArr(param string) []int64 {
	params := strings.Split(strings.TrimSpace(param), ",")
	var result []int64

	for _, p := range params {
		i, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			logging.Debug.Printf("[queryParamToIntArr] Invalid params: %+v", params)
		} else {
			result = append(result, i)
		}
	}

	return result
}

func ValidateStatus(ctx context.Context, status []int64) bool {
	for i := 0; i < len(status); i++ {
		if status[i] < 0 {
			return false
		}
	}
	return true
}

func StringArrayToMap(ctx context.Context, ar []string) map[string]bool {
	result := make(map[string]bool)
	for i := range ar {
		result[ar[i]] = true
	}
	return result
}

func RawstringToStructure(ctx context.Context, raw string, structData interface{}) error {

	rawIn := json.RawMessage(raw)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		logging.Debug.Println("[RawstringToStructure] Unable to Marshal input. Err: ", err)
		return err
	}

	err = json.Unmarshal(bytes, structData)
	if err != nil {
		logging.Debug.Println("[RawstringToStructure] Unable to UnMarshal input. Err: ", err)
		return err
	}

	return nil
}

func TimeByZone(ctx context.Context, location, format string, tx time.Time) string {
	var defaultString string

	if format == "" {
		format = constant.DateFormat
	}

	if location == "" {
		location = constant.LOC_JAKARTA
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		return defaultString
	}

	now := tx.In(loc).Format(format)

	return now
}

func RemoveDuplicates(arr []int64) (result []int64) {
	list := make(map[int64]bool)

	for _, val := range arr {
		list[val] = true
	}

	for key := range list {
		result = append(result, key)
	}

	return result
}

func IntArrayToMap(ctx context.Context, list []int64) map[int64]bool {
	result := make(map[int64]bool)
	for _, val := range list {
		result[val] = true
	}
	return result
}

// COnverts the given time to the time according to the location provided
func ConvertTimeToTimeZone(ctx context.Context, currentTime time.Time, location string, offset int) (*time.Time, error) {

	var (
		loc         *time.Location
		updatedTime time.Time
	)

	loc = time.FixedZone(location, offset)
	updatedTime = currentTime.In(loc)
	return &updatedTime, nil
}

func TruncateString(ctx context.Context, inputString string, maxLength int) string {
	if len(inputString) > maxLength {
		inputString = inputString[:maxLength]
	}
	return inputString
}

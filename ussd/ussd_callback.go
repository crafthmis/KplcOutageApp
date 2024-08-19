package ussd

import (
	"encoding/json"
	"fmt"
	"kplc-outage-app/models"
	"kplc-outage-app/services"
	"net/http"
	"strconv"
	"strings"
)

var lastIndexOf int

func UssdCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	// recieve formValue from AT
	session_id := r.FormValue("sessionId")
	service_code := r.FormValue("serviceCode")
	phone_number := r.FormValue("phoneNumber")
	text := r.FormValue("text")

	fmt.Printf("%s,%s,%s", session_id, service_code, phone_number)

	// if the text field is empty, this indicates that this is the begining of a session
	if len(text) == 0 {
		// form the response to be sent back to the user
		var Subscriptions []models.Subscription

		err := services.GetAllSubscriptions(&Subscriptions)
		if err != nil {
			w.Write([]byte("END System is currently busy. Kindly try again"))
			return
		}
		var accumulator []string

		for _, sub := range Subscriptions {
			result := fmt.Sprintf("\n%d. %s", sub.SubID, sub.Name)
			accumulator = append(accumulator, result)
		}

		err = services.CreateUssdSession(&models.UssdSession{SessionID: session_id, Msisdn: phone_number})
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte("END System is currently busy. Kindly try again"))
			return
		}
		output := strings.Join(accumulator, "")
		fmt.Println(output)
		w.Write([]byte(fmt.Sprintf("CON Welcome to nijulishe. Choose you plan. %s", output)))
		return
	} else {
		//   On user input the switch block is executed, remember our text field is concatenated on every user input
		cnt := strings.Count(text, "*")
		fmt.Printf("\nHow many astericks %d", cnt)
		lastIndexOf = getLastValueAfterAsterisk(text)
		fmt.Printf("\nLast Index %d", lastIndexOf)
		switch cnt {
		case 0:
			var regions []models.Region
			err := services.GetAllRegion(&regions)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			var accumulator []string
			for idx, reg := range regions {
				result := fmt.Sprintf("\n%d. %s", idx+1, reg.Name)
				accumulator = append(accumulator, result)
			}
			jsonData, err := json.MarshalIndent(regions, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling to JSON:", err)
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			updates := map[string]interface{}{
				"region_payload": string(jsonData),
			}
			err = services.UpdateUssdSession(updates, session_id)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}

			output := strings.Join(accumulator, "")
			fmt.Println(output)
			w.Write([]byte(fmt.Sprintf("CON Choose your Region. %s", output)))
			return
		case 1:
			var region models.Region
			var session models.UssdSession
			err := services.GetUssdSessionByID(&session, session_id)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}

			var regions []models.Region
			err = json.Unmarshal([]byte(session.RegionPayload), &regions)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				fmt.Println("Error unmarshalling JSON:", err)
				return
			}

			reg, err := GetRegionByIndex(regions, lastIndexOf-1)

			fmt.Printf("\nRegion Name %s", reg.Name)

			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			err = services.GetRegionByID(&region, fmt.Sprintf("%d", reg.RegID))
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			var accumulator []string
			for idx, county := range region.Counties {
				result := fmt.Sprintf("\n%d. %s", idx+1, county.Name)
				accumulator = append(accumulator, result)
			}
			// Save counties
			jsonData, err := json.MarshalIndent(region.Counties, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling to JSON:", err)
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			updates := map[string]interface{}{
				"county_payload": string(jsonData),
			}
			err = services.UpdateUssdSession(updates, session_id)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}

			output := strings.Join(accumulator, "")
			fmt.Println(output)
			w.Write([]byte(fmt.Sprintf("CON Choose your County. %s", output)))
			return
		case 2:
			var county models.County
			var session models.UssdSession
			err := services.GetUssdSessionByID(&session, session_id)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}

			var counties []models.County
			err = json.Unmarshal([]byte(session.CountyPayload), &counties)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				fmt.Println("Error unmarshalling JSON:", err)
				return
			}

			cout, err := GetCountyByIndex(counties, lastIndexOf-1)

			fmt.Printf("\nCounty Name %s", cout.Name)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			err = services.GetCountyByID(&county, fmt.Sprintf("%d", cout.CtyID))
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}

			var accumulator []string
			for idx, constituency := range county.Constituencies {
				result := fmt.Sprintf("\n%d. %s", idx+1, constituency.Name)
				accumulator = append(accumulator, result)
			}
			// Save constituencies
			jsonData, err := json.MarshalIndent(county.Constituencies, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling to JSON:", err)
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			fmt.Printf("\njsonData %s", jsonData)
			updates := map[string]interface{}{
				"constituency_payload": string(jsonData),
			}
			err = services.UpdateUssdSession(updates, session_id)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			output := strings.Join(accumulator, "")
			fmt.Println(output)
			w.Write([]byte(fmt.Sprintf("CON Choose your Constituency. %s", output)))
			return
		case 3:
			var constituency models.Constituency
			var session models.UssdSession
			err := services.GetUssdSessionByID(&session, session_id)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}

			var constituencies []models.Constituency
			err = json.Unmarshal([]byte(session.ConstituencyPayload), &constituencies)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			fmt.Printf("\nConstituencyPayload %s", session.ConstituencyPayload)

			con, err := GetConstituencyByIndex(constituencies, lastIndexOf-1)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			fmt.Printf("\nConstituency Name %s", con.Name)

			err = services.GetConstituencyByID(&constituency, fmt.Sprintf("%d", con.CstID))
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			var accumulator []string
			for idx, area := range constituency.Areas {
				result := fmt.Sprintf("\n%d. %s", idx+1, area.Name)
				accumulator = append(accumulator, result)
			}

			// Save areas
			jsonData, err := json.MarshalIndent(constituency.Areas, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling to JSON:", err)
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			fmt.Printf("\njsonData %s", jsonData)
			updates := map[string]interface{}{
				"area_payload": string(jsonData),
			}
			err = services.UpdateUssdSession(updates, session_id)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			output := strings.Join(accumulator, "")
			fmt.Println(output)
			w.Write([]byte(fmt.Sprintf("CON Choose your Area. %s", output)))
			return
		case 4:
			var session models.UssdSession
			err := services.GetUssdSessionByID(&session, session_id)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}

			var areas []models.Area
			err = json.Unmarshal([]byte(session.AreaPayload), &areas)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			fmt.Printf("\nConstituencyPayload %s", session.ConstituencyPayload)

			are, err := GetAreaByIndex(areas, lastIndexOf-1)
			if err != nil {
				w.Write([]byte("END System is currently busy. Kindly try again"))
				return
			}
			fmt.Printf("\nArea Name %s", are.Name)

			w.Write([]byte("END You have been registered successfully. Thank you."))
			return
		default:
			w.Write([]byte("END Invalid input"))
			return
		}
	}

}

func getLastValueAfterAsterisk(s string) int {
	if !strings.Contains(s, "*") {
		return -1 // Return the original string if no asterisk is found
	}
	// Split the string by asterisks
	parts := strings.Split(s, "*")

	// Check if there are any parts after splitting
	if len(parts) == 0 {
		return -1
	}

	str := strings.TrimSpace(parts[len(parts)-1])
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return -1
	}
	// Return the last part (trimmed of any whitespace)
	return num
}

func GetRegionByIndex(regions []models.Region, idx int) (models.Region, error) {

	if idx < 0 || idx >= len(regions) {
		return models.Region{}, fmt.Errorf("index out of range")
	}
	return regions[idx], nil
}

func GetCountyByIndex(counties []models.County, idx int) (models.County, error) {

	if idx < 0 || idx >= len(counties) {
		return models.County{}, fmt.Errorf("index out of range")
	}
	return counties[idx], nil
}

func GetConstituencyByIndex(constituencies []models.Constituency, idx int) (models.Constituency, error) {

	if idx < 0 || idx >= len(constituencies) {
		return models.Constituency{}, fmt.Errorf("index out of range")
	}
	return constituencies[idx], nil
}

func GetAreaByIndex(areas []models.Area, idx int) (models.Area, error) {

	if idx < 0 || idx >= len(areas) {
		return models.Area{}, fmt.Errorf("index out of range")
	}
	return areas[idx], nil
}

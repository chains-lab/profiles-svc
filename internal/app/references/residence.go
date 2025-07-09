package references

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func ValidateResidence(city, region, country, apiKey string) error {
	countryCode, err := getCountryCode(country, apiKey)
	if err != nil {
		return fmt.Errorf("invalid country: %w", err)
	}

	regionCode, err := getRegionCode(region, countryCode, apiKey)
	if err != nil {
		return fmt.Errorf("invalid region for country '%s': %w", country, err)
	}

	found, err := cityInRegion(city, regionCode, apiKey)
	if err != nil {
		return fmt.Errorf("failed to validate city: %w", err)
	}
	if !found {
		return fmt.Errorf("city '%s' not found in region '%s' of country '%s'", city, region, country)
	}

	return nil
}

func getCountryCode(name, apiKey string) (string, error) {
	url := "https://wft-geo-db.p.rapidapi.com/v1/geo/countries?namePrefix=" + url.QueryEscape(name)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", "wft-geo-db.p.rapidapi.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	for _, c := range result.Data {
		if strings.EqualFold(c.Name, name) {
			return c.Code, nil
		}
	}
	return "", fmt.Errorf("country not found")
}

func getRegionCode(regionName, countryCode, apiKey string) (string, error) {
	url := fmt.Sprintf("https://wft-geo-db.p.rapidapi.com/v1/geo/countries/%s/regions?namePrefix=%s", countryCode, url.QueryEscape(regionName))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", "wft-geo-db.p.rapidapi.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	for _, r := range result.Data {
		if strings.EqualFold(r.Name, regionName) {
			return r.Code, nil
		}
	}
	return "", fmt.Errorf("region not found")
}

func cityInRegion(city, regionCode, apiKey string) (bool, error) {
	url := fmt.Sprintf("https://wft-geo-db.p.rapidapi.com/v1/geo/regions/%s/cities?namePrefix=%s", regionCode, url.QueryEscape(city))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", "wft-geo-db.p.rapidapi.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	for _, c := range result.Data {
		if strings.EqualFold(c.Name, city) {
			return true, nil
		}
	}
	return false, nil
}

package main

import (
	"fmt"
	"github.com/icrowley/fake"

	"golang.org/x/net/context"
	"cloud.google.com/go/bigquery"
	"log"
	"github.com/matejl/challenge/model"
	"github.com/Pallinder/go-randomdata"
	"time"
	"github.com/matejl/challenge/database"
)

const (

	//MinCampaigns = 3
	//MaxCampaigns = 4
	//
	//MinAds = 1
	//MaxAds = 50
	//
	//MinUsers = 1
	//MaxUsers = 500000
	//
	//MinImpressions = 1
	//MaxImpressions = 10

	MinCampaigns = 2
	MaxCampaigns = 4

	MinAds = 2
	MaxAds = 5

	MinUsers = 5
	MaxUsers = 10

	MinImpressions = 1
	MaxImpressions = 4

	// 2 years
	MaxDateDaysBack = 730

	MaxClicks       = 2
	MaxInteractions = 2
)

// POST https://www.googleapis.com/bigquery/v2/projects/projectId/datasets/datasetId/tables/tableId/insertAll

func main() {
	// Print a random silly name
	fmt.Println(fake.Company() + " " + fake.ProductName())

	client, err := database.GetStat()
	if err != nil {
		log.Fatal(err)
	}

	dataset := client.Dataset(model.DatasetId)
	//log.Printf("%#v", getData())

	saveData := NewSaveData()

	log.Println("Table")
	campaignTable := dataset.Table("campaign")
	adTable := dataset.Table("ad")
	userTable := dataset.Table("user")
	impressionTable := dataset.Table("impression")

	log.Println("Putting campaigns")
	campaignUploader := campaignTable.Uploader()
	err = campaignUploader.Put(context.Background(), saveData.campaigns)
	if err != nil {
		log.Fatal("Error when inserting campaigns:", err)
	}

	log.Println("Putting ads")
	adUploader := adTable.Uploader()
	err = adUploader.Put(context.Background(), saveData.ads)
	if err != nil {
		log.Fatal("Error when inserting ads:", err)
	}

	log.Println("Putting users")
	userUploader := userTable.Uploader()
	err = userUploader.Put(context.Background(), saveData.users)
	if err != nil {
		log.Fatal("Error when inserting users:", err)
	}

	log.Println("Putting impressions")
	impressionUploader := impressionTable.Uploader()
	err = impressionUploader.Put(context.Background(), saveData.impressions)
	if err != nil {
		multiError := err.(bigquery.PutMultiError)
		for _, e := range multiError {
			log.Println(e)
		}
		log.Fatal("Error when inserting impressions:", err)
	}

	log.Println("Done")

}

type SaveData struct {
	campaigns   []*model.Campaign
	ads         []*model.Ad
	users       []*model.User
	impressions []*model.Impression

	campaignId   int64
	adId         int64
	userId       int64
	impressionId int64
}

func NewSaveData() SaveData {

	sd := SaveData{ }
	if maxCampaignId, err := model.MaxCampaignId(); err == nil {
		sd.campaignId = maxCampaignId + 1
	}
	if maxAdId, err := model.MaxAdId(); err == nil {
		sd.adId = maxAdId + 1
	}
	if maxUserId, err := model.MaxUserId(); err == nil {
		sd.userId = maxUserId + 1
	}
	if maxImpressionId, err := model.MaxImpressionId(); err == nil {
		sd.impressionId = maxImpressionId + 1
	}

	sd.campaigns = make([]*model.Campaign, 0)
	sd.ads = make([]*model.Ad, 0)
	sd.users = make([]*model.User, 0)
	sd.impressions = make([]*model.Impression, 0)

	nCampaigns := int64(randomdata.Number(MinCampaigns, MaxCampaigns))

	for ; nCampaigns > 0; nCampaigns-- {
		log.Println("Creating campaign id ", sd.campaignId)
		sd.campaigns = append(sd.campaigns, sd.newCampaign(sd.campaignId))

		nAds := int64(randomdata.Number(MinAds, MaxAds))
		log.Println("nAds: ", nAds)

		for ; nAds > 0; nAds-- {
			sd.ads = append(sd.ads, sd.newAd(sd.adId, sd.campaignId))

			nUsers := int64(randomdata.Number(MinUsers, MaxUsers))
			log.Println("nUsers: ", nUsers)
			for ; nUsers > 0; nUsers-- {

				if sd.userId%10000 == 0 {
					log.Println("User id ", sd.userId)
				}

				sd.users = append(sd.users, sd.newUser(sd.userId, sd.adId))

				nImpressions := int64(randomdata.Number(MinImpressions, MaxImpressions))
				for ; nImpressions > 0; nImpressions-- {
					sd.impressions = append(sd.impressions, sd.newImpression(sd.impressionId, sd.userId))
					sd.impressionId++
				}
				sd.userId++
			}
			sd.adId++
		}
		sd.campaignId++
	}

	return sd

}

func (sd *SaveData) newCampaign(campaignId int64) *model.Campaign {

	return &model.Campaign{
		CampaignId:   campaignId,
		CampaignName: fake.Company() + " " + fake.ProductName(),
	}

}

func (sd *SaveData) newAd(adId int64, campaignId int64) *model.Ad {

	return &model.Ad{
		AdId:       adId,
		CampaignId: campaignId,
		AdName:     fake.Word(),
	}

}

func (sd *SaveData) newUser(userId int64, adId int64) *model.User {

	return &model.User{
		UserId:   userId,
		AdId:     adId,
		UserHash: randomdata.Letters(16),
	}

}

func (sd *SaveData) newImpression(impressionId int64, userId int64) *model.Impression {

	return &model.Impression{
		ImpressionId: impressionId,
		UserId:       userId,
		Click:        randomdata.Number(MaxClicks),
		Datetime:     time.Now().AddDate(0, 0, -1*randomdata.Number(MaxDateDaysBack)),
		Pinch:        randomdata.Number(MaxInteractions),
		Swipe:        randomdata.Number(MaxInteractions),
		Touch:        randomdata.Number(MaxInteractions),
	}

}

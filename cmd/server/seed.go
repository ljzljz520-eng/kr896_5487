package main

import (
	"ruralfolk/domain"
	"ruralfolk/service"
	"ruralfolk/store"
)

func seedFixtures(s *store.Store) error {
	svc := service.New(s)
	fixtures := []domain.Exhibit{{ID: "fixture-1", Title: "竹编灯笼", Story: "夜里点亮村口的手艺", Status: domain.Published, MediaURL: "/media/lantern.jpg"}, {ID: "fixture-2", Title: "土布织机", Story: "经纬之间保存着节气", Status: domain.Published, MediaURL: "/media/loom.jpg"}, {ID: "fixture-3", Title: "木版年画", Story: "一笔一画都是家门口的祝愿", Status: domain.Published, MediaURL: "/media/print.jpg"}}
	for _, e := range fixtures {
		if err := s.SaveExhibit(e); err != nil {
			return err
		}
	}
	artisans := []domain.Artisan{{ID: "artisan-1", Name: "赵桂兰", Bio: "四十年只做一件布衣", Craft: "土布织造", PortraitURL: "/media/zhao.jpg"}, {ID: "artisan-2", Name: "周大山", Bio: "把竹子编成会呼吸的灯", Craft: "竹编", PortraitURL: "/media/zhou.jpg"}}
	for _, a := range artisans {
		if err := svc.SaveArtisan(a); err != nil {
			return err
		}
	}
	news := []domain.News{{ID: "news-1", Title: "秋收后的织布课", Body: "在老粮仓认识经纬", Published: true}, {ID: "news-2", Title: "木版印色工作坊", Body: "跟着手艺人刻一张年画", Published: true}}
	for _, n := range news {
		if err := svc.SaveNews(n); err != nil {
			return err
		}
	}
	return nil
}
func fixtureExists(s *store.Store, id string) bool {
	e, err := s.GetExhibit(id)
	return err == nil && e.ID == id
}
func ensureFixtures(s *store.Store) error {
	if fixtureExists(s, "fixture-1") {
		return nil
	}
	return seedFixtures(s)
}

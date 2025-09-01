package _const

// ProductTypeMap 一级分类映射
var ProductTypeMap = map[int64]string{
	ProductTypeDigital:              "数码",
	ProductTypeMusicalinstruments:   "玩具乐器",
	ProductTypeHomeappliances:       "家电",
	ProductTypeMedicine:             "医药",
	ProductTypePetflowers:           "宠物鲜花",
	ProductTypeWomenfashion:         "女装",
	ProductTypeOutdoorSports:        "户外运动",
	ProductTypeluxurygoods:          "奢侈品",
	ProductTypeUsed:                 "二手",
	ProductTypeFoodbeverages:        "食品饮料",
	ProductTypeUnderwearaccessories: "内衣配饰",
	ProductTypeMobilephone:          "手机",
	ProductTypeBoot:                 "靴子",
	ProductTypeBookEntertainment:    "图书文娱",
	ProductTypeClockGlasses:         "钟表眼镜",
	ProductTypeIndustrialproducts:   "工业品",
	ProductTypeMkw:                  "母婴童装",
	ProductTypeGoldJewelry:          "黄金珠宝",
	ProductTypeFurniture:            "家具家装",
	ProductTypeKitchenware:          "厨具",
	ProductTypeComponents:           "元器件",
	ProductTypeStationery:           "文具",
	ProductTypeLuggage:              "箱包",
	ProductTypeMenswear:             "男装",
	ProductTypeBeautyskin:           "美容护肤",
	ProductTypeHousehold:            "家居",
	ProductTypeOfficeworkcomputers:  "电脑办公",
}

// ProductSubTypeMap 二级分类映射
var ProductSubTypeMap = map[int64]string{
	// 数码二级分类
	ProductTypeDigitalSmartTerminals:      "超级趋势",
	ProductTypeDigitalSurveillance:        "智能监控",
	ProductTypeDigitalSmartdevice:         "智能设备",
	ProductTypeDigitalCamera:              "摄像机",
	ProductTypeDigitalAccessories:         "配件",
	ProductTypeDigitalFilmentertainment:   "影视娱乐",
	ProductTypeDigitalElectroniceducation: "电子教育",

	// 玩具乐器二级分类
	ProductTypeMusicalinstrumentsTrending:             "热门",
	ProductTypeMusicalinstrumentsFashionableanime:     "潮玩动漫",
	ProductTypeMusicalinstrumentsCreativestressrelief: "创意减压",
	ProductTypeMusicalinstrumentsEarlyeducation:       "早教益智",
	ProductTypeMusicalinstrumentsToysandgames:         "玩具游戏",
	ProductTypeMusicalinstrumentsMusicalinstruments:   "乐器",

	// 家电二级分类
	ProductTypeHomeappliancesTrending:        "热门",
	ProductTypeHomeappliancesSmartappliances: "智能家居",
	ProductTypeHomeapplianceswatch:           "手表",
	ProductTypeHomeapplianceslight:           "灯",
	ProductTypeHomeapplianceslock:            "门锁",
	ProductTypeHomeappliancespeaker:          "音箱",
	ProductTypeHomeappliancesTV:              "电视",
	ProductTypeHomeappliancesefrigerator:     "冰箱",
	ProductTypeHomeapplianceswashingmachine:  "洗衣机",
	ProductTypeHomeappliancestairconditioner: "空调",

	// 医药二级分类
	ProductTypeMedicineTrending:    "热门",
	ProductTypeMedicinehealthcare:  "健康医疗",
	ProductTypeMedicinedrugstore:   "药店",
	ProductTypeMedicinedrug:        "药品",
	ProductTypeMedicinenutrition:   "营养保健",
	ProductTypeMedicinenursinghome: "护理",

	// 宠物鲜花二级分类
	ProductTypePetflowersTrending: "热门",
	ProductTypePetflowerspetfood:  "宠物食品",
	ProductTypePetflowerspettoys:  "宠物玩具",
	ProductTypePetflowerfresh:     "鲜花",

	// 女装二级分类
	ProductTypeWomenfashionTrending: "热门",
	ProductTypeWomenfashionclothing: "女装服装",
	ProductTypeWomenfashionshoes:    "女装鞋",
	ProductTypeWomenfashionjewelry:  "女装饰品",
	ProductTypeWomenfashionbags:     "女装包",

	// 户外运动二级分类
	ProductTypeOutdoorSportsTrending: "热门",
	ProductTypeOutdoorSportsbiking:   "自行车",
	ProductTypeOutdoorSportsskiing:   "滑雪",
	ProductTypeOutdoorSportsswimming: "游泳",
	ProductTypeOutdoorSportscamping:  "露营",
	ProductTypeOutdoorSportstennis:   "网球",

	// 奢侈品二级分类
	ProductTypeLuxurygoodsTrending: "热门",
	ProductTypeLuxurygoodsjewelry:  "珠宝",
	ProductTypeLuxurygoodsearings:  "首饰",
	ProductTypeLuxurygoodscar:      "汽车",
	ProductTypeLuxurygoodsshoes:    "鞋",
	ProductTypeLuxurygoodscarpets:  "皮包",

	// 二手二级分类
	ProductTypeUsedTrending:    "热门",
	ProductTypeUsedclothing:    "服装",
	ProductTypeUsedshoes:       "鞋",
	ProductTypeUsedjewelry:     "饰品",
	ProductTypeUsedbags:        "包",
	ProductTypeUsedmobilephone: "手机",
	ProductTypeUsedcomputer:    "电脑",
	ProductTypeUsedbooks:       "书",
	ProductTypeUsedstationery:  "文具",
	ProductTypeUsedfurniture:   "家具",
	ProductTypeUsedkitchenware: "厨具",

	// 食品饮料二级分类
	ProductTypeFoodbeveragesTrending:     "热门",
	ProductTypeFoodbeveragesgroceries:    "生鲜",
	ProductTypeFoodbeveragesmeat:         "肉类",
	ProductTypeFoodbeveragesvegetables:   "蔬菜",
	ProductTypeFoodbeveragesfruits:       "水果",
	ProductTypeFoodbeveragesseafood:      "海鲜",
	ProductTypeFoodbeveragesgrain:        "谷物",
	ProductTypeFoodbeveragescoffee:       "咖啡",
	ProductTypeFoodbeveragesenergydrinks: "能量饮料",
	ProductTypeFoodbeveragesjuices:       "果汁",

	// 内衣配饰二级分类
	ProductTypeUnderwearaccessoriesTrending: "热门",
	ProductTypeUnderwearaccessoriesclothing: "女士内衣",
	ProductTypeUnderwearaccessoriesmanshoes: "男士内衣",

	// 手机二级分类
	ProductTypeMobilephoneTrending:    "热门",
	ProductTypeMobilephonesmart:       "智能手机",
	ProductTypeMobilephonephone:       "手机壳",
	ProductTypeMobilephoneaccessories: "配件",
	ProductTypeMobilephonecharger:     "充电器",

	// 靴子二级分类
	ProductTypeBootTrending:  "热门",
	ProductTypeBootsfootwear: "靴子",

	// 图书文娱二级分类
	ProductTypeBookEntertainmentTrending:  "热门",
	ProductTypeBookEntertainmentnovels:    "小说",
	ProductTypeBookEntertainmentcomics:    "漫画",
	ProductTypeBookEntertainmentmagazines: "杂志",
	ProductTypeBookEntertainmentmusic:     "音乐",
	ProductTypeBookEntertainmentmovies:    "电影",

	// 钟表眼镜二级分类
	ProductTypeClockGlassesTrending: "热门",
	ProductTypeClockGlasseswatches:  "表带",
	ProductTypeClockGlassesclocks:   "手表",
	ProductTypeClockGlassesglasses:  "眼镜",

	// 工业品二级分类
	ProductTypeIndustrialproductsTrending:     "热门",
	ProductTypeIndustrialproductstools:        "工具",
	ProductTypeIndustrialproductscables:       "清洁用品",
	ProductTypeIndustrialproductscontrolpanel: "控制器",

	// 母婴童装二级分类
	ProductTypeMkwTrending:     "热门",
	ProductTypeMkwbabyproducts: "婴童用品",
	ProductTypeMkwtoys:         "玩具",
	ProductTypeMkwbabydiapers:  "奶粉",
	ProductTypeMkwclothing:     "服装",
	ProductTypeMkwbabyshoes:    "童鞋",

	// 黄金珠宝二级分类
	ProductTypeGoldJewelryTrending:  "热门",
	ProductTypeGoldJewelrygold:      "黄金",
	ProductTypeGoldJewelryjewelry:   "珠宝",
	ProductTypeGoldJewelryrings:     "戒指",
	ProductTypeGoldJewelrynecklaces: "项链",

	// 家具家装二级分类
	ProductTypeFurnitureTrending: "热门",
	ProductTypeFurnitureschair:   "客厅",
	ProductTypeFurnitureshelves:  "书架",
	ProductTypeFurnituretable:    "柜子",

	// 厨具二级分类
	ProductTypeKitchenwareTrending:          "热门",
	ProductTypeKitchenwarecutlery:           "刀剪菜板",
	ProductTypeKitchenwarecookware:          "厨房用具",
	ProductTypeKitchenwareutensils:          "餐具",
	ProductTypeKitchenwarecookingutensils:   "锅",
	ProductTypeKitchenwarekitchenappliances: "厨房电器",

	// 元器件二级分类
	ProductTypeComponentsTrending:     "热门",
	ProductTypeComponentscables:       "电缆",
	ProductTypeComponentssensors:      "传感器",
	ProductTypeComponentscontrolpanel: "控制器",
	ProductTypeComponentscontrolunit:  "控制单元",
	ProductTypeComponentscontrolboard: "控制板",

	// 文具二级分类
	ProductTypeStationeryTrending: "热门",
	ProductTypeStationerypen:      "笔",
	ProductTypeStationerypencil:   "铅笔",
	ProductTypeStationeryruler:    "尺子",
	ProductTypeStationeryeraser:   "橡皮",
	ProductTypeStationerypaper:    "纸",
	ProductTypeStationeryink:      "墨水",

	// 箱包二级分类
	ProductTypeLuggageTrending:  "热门",
	ProductTypeLuggagebags:      "包包",
	ProductTypeLuggagebackpacks: "手提包",
	ProductTypeLuggagewallets:   "钱包",
	ProductTypeLuggagesuitcases: "行李箱",

	// 男装二级分类
	ProductTypeMenswearTrending: "热门",
	ProductTypeMenswearclothing: "男装服装",
	ProductTypeMenswearshoes:    "男装鞋",
	ProductTypeMenswearjewelry:  "男装饰品",
	ProductTypeMenswearbags:     "男包",

	// 美容护肤二级分类
	ProductTypeBeautyskinTrending: "热门",
	ProductTypeBeautyskinface:     "面部护肤",
	ProductTypeBeautyskinbody:     "身体护肤",
	ProductTypeBeautyskinmakeup:   "美妆",
	ProductTypeBeautyskinskincare: "皮肤护理",

	// 家居二级分类
	ProductTypeHouseholdTrending:   "热门",
	ProductTypeHouseholdfurniture:  "家具",
	ProductTypeHouseholddecorative: "装饰",
	ProductTypeHouseholdlighting:   "照明",
	ProductTypeHouseholdappliances: "家电",

	// 电脑办公二级分类
	ProductTypeOfficeworkcomputersTrending:    "热门",
	ProductTypeOfficeworkcomputersdesktop:     "台式机",
	ProductTypeOfficeworkcomputerslaptop:      "笔记本",
	ProductTypeOfficeworkcomputersprinter:     "打印机",
	ProductTypeOfficeworkcomputersaccessories: "配件",
	ProductTypeOfficeworkcomputerssoftware:    "软件",
	ProductTypeOfficeworkcomputersnetwork:     "网络设备",
}

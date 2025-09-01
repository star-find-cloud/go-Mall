package _const

const (
	//一级分类（大类）
	ProductTypeNil                  = iota        //无分类
	ProductTypeDigital              = iota*10 + 1 //数码
	ProductTypeMusicalinstruments                 //玩具乐器
	ProductTypeHomeappliances                     //家电
	ProductTypeMedicine                           //医药
	ProductTypePetflowers                         //宠物鲜花
	ProductTypeWomenfashion                       //女装
	ProductTypeOutdoorSports                      //户外运动
	ProductTypeluxurygoods                        //奢侈品
	ProductTypeUsed                               //二手
	ProductTypeFoodbeverages                      //食品饮料
	ProductTypeUnderwearaccessories               //内衣配饰
	ProductTypeMobilephone                        //手机
	ProductTypeBoot                               //靴子
	ProductTypeBookEntertainment                  //图书文娱
	ProductTypeClockGlasses                       //钟表眼镜
	ProductTypeIndustrialproducts                 //工业品
	ProductTypeMkw                                //母婴童装
	ProductTypeGoldJewelry                        //黄金珠宝
	ProductTypeFurniture                          //家具家装
	ProductTypeKitchenware                        //厨具
	ProductTypeComponents                         //元器件
	ProductTypeStationery                         //文具
	ProductTypeLuggage                            //箱包
	ProductTypeMenswear                           //男装
	ProductTypeBeautyskin                         //美容护肤
	ProductTypeHousehold                          //家居
	ProductTypeOfficeworkcomputers                //电脑办公
)

const (
	//二级分类（数码）
	ProductTypeDigitalSmartTerminals      = (ProductTypeDigital * 100) + iota + 1 //智能终端101
	ProductTypeDigitalSurveillance                                                //智能监控102
	ProductTypeDigitalSmartdevice                                                 //智能设备103
	ProductTypeDigitalCamera                                                      //摄像机104
	ProductTypeDigitalAccessories                                                 //配件105
	ProductTypeDigitalFilmentertainment                                           //影视娱乐106
	ProductTypeDigitalElectroniceducation                                         //电子教育107
)

const (
	//二级分类（玩具乐器）
	ProductTypeMusicalinstrumentsTrending             = (ProductTypeMusicalinstruments * 100) + iota + 1 //热门201
	ProductTypeMusicalinstrumentsFashionableanime                                                        //潮玩动漫202
	ProductTypeMusicalinstrumentsCreativestressrelief                                                    //创意减压203
	ProductTypeMusicalinstrumentsEarlyeducation                                                          //早教益智204
	ProductTypeMusicalinstrumentsToysandgames                                                            //玩具游戏205
	ProductTypeMusicalinstrumentsMusicalinstruments                                                      //乐器206
)

const (
	//二级分类（家电）
	ProductTypeHomeappliancesTrending        = (ProductTypeHomeappliances * 100) + iota + 1 //热门301
	ProductTypeHomeappliancesSmartappliances                                                //智能家居302
	ProductTypeHomeapplianceswatch                                                          //手表303
	ProductTypeHomeapplianceslight                                                          //灯304
	ProductTypeHomeapplianceslock                                                           //门锁305
	ProductTypeHomeappliancespeaker                                                         //音箱306
	ProductTypeHomeappliancesTV                                                             //电视307
	ProductTypeHomeappliancesefrigerator                                                    //冰箱308
	ProductTypeHomeapplianceswashingmachine                                                 //洗衣机309
	ProductTypeHomeappliancestairconditioner                                                //空调310
)

const (
	//二级分类（医药）
	ProductTypeMedicineTrending    = (ProductTypeMedicine * 100) + iota + 1 //热门401
	ProductTypeMedicinehealthcare                                           //健康医疗402
	ProductTypeMedicinedrugstore                                            //药店403
	ProductTypeMedicinedrug                                                 //药品404
	ProductTypeMedicinenutrition                                            //营养保健405
	ProductTypeMedicinenursinghome                                          //护理406
)

const (
	//二级分类（宠物鲜花）
	ProductTypePetflowersTrending = (ProductTypePetflowers * 100) + iota + 1 //热门501
	ProductTypePetflowerspetfood                                             //宠物食品502
	ProductTypePetflowerspettoys                                             //宠物玩具503
	ProductTypePetflowerfresh                                                //鲜花504
)

const (
	//二级分类（女装）
	ProductTypeWomenfashionTrending = (ProductTypeWomenfashion * 100) + iota + 1 //热门601
	ProductTypeWomenfashionclothing                                              //女装服装602
	ProductTypeWomenfashionshoes                                                 //女装鞋603
	ProductTypeWomenfashionjewelry                                               //女装饰品604
	ProductTypeWomenfashionbags                                                  //女装包605
)

const (
	//二级分类（户外运动）
	ProductTypeOutdoorSportsTrending = (ProductTypeOutdoorSports * 100) + iota + 1 //热门701
	ProductTypeOutdoorSportsbiking                                                 //自行车702
	ProductTypeOutdoorSportsskiing                                                 //滑雪703
	ProductTypeOutdoorSportsswimming                                               //游泳704
	ProductTypeOutdoorSportscamping                                                //露营705
	ProductTypeOutdoorSportstennis                                                 //网球706
)

const (
	//二级分类（奢侈品）
	ProductTypeLuxurygoodsTrending = (ProductTypeluxurygoods * 100) + iota + 1 //热门801
	ProductTypeLuxurygoodsjewelry                                              //珠宝802
	ProductTypeLuxurygoodsearings                                              //首饰803
	ProductTypeLuxurygoodscar                                                  //汽车804
	ProductTypeLuxurygoodsshoes                                                //鞋805
	ProductTypeLuxurygoodscarpets                                              //皮包806
)

const (
	//二级分类（二手）
	ProductTypeUsedTrending    = (ProductTypeUsed * 100) + iota + 1 //热门901
	ProductTypeUsedclothing                                         //服装902
	ProductTypeUsedshoes                                            //鞋903
	ProductTypeUsedjewelry                                          //饰品904
	ProductTypeUsedbags                                             //包905
	ProductTypeUsedmobilephone                                      //手机906
	ProductTypeUsedcomputer                                         //电脑907
	ProductTypeUsedbooks                                            //书908
	ProductTypeUsedstationery                                       //文具909
	ProductTypeUsedfurniture                                        //家具910
	ProductTypeUsedkitchenware                                      //厨具911
)

const (
	//二级分类（食品饮料）
	ProductTypeFoodbeveragesTrending     = (ProductTypeFoodbeverages * 100) + iota + 1 //热门1001
	ProductTypeFoodbeveragesgroceries                                                  //生鲜1002
	ProductTypeFoodbeveragesmeat                                                       //肉类1003
	ProductTypeFoodbeveragesvegetables                                                 //蔬菜1004
	ProductTypeFoodbeveragesfruits                                                     //水果1005
	ProductTypeFoodbeveragesseafood                                                    //海鲜1006
	ProductTypeFoodbeveragesgrain                                                      //谷物1007
	ProductTypeFoodbeveragescoffee                                                     //咖啡1008
	ProductTypeFoodbeveragesenergydrinks                                               //能量饮料1009
	ProductTypeFoodbeveragesjuices                                                     //果汁1010
)

const (
	//二级分类（内衣配饰）
	ProductTypeUnderwearaccessoriesTrending = (ProductTypeUnderwearaccessories * 100) + iota + 1 //热门1101
	ProductTypeUnderwearaccessoriesclothing                                                      //女士内衣1102
	ProductTypeUnderwearaccessoriesmanshoes                                                      //男士内衣1102
)

const (
	//二级分类（手机）
	ProductTypeMobilephoneTrending    = (ProductTypeMobilephone * 100) + iota + 1 //热门1201
	ProductTypeMobilephonesmart                                                   //智能手机1202
	ProductTypeMobilephonephone                                                   //手机壳1203
	ProductTypeMobilephoneaccessories                                             //配件1204
	ProductTypeMobilephonecharger                                                 //充电器1205
)

const (
	//二级分类（靴子）
	ProductTypeBootTrending  = (ProductTypeBoot * 100) + iota + 1 //热门1301
	ProductTypeBootsfootwear                                      //靴子1302
)

const (
	//二级分类（图书文娱）
	ProductTypeBookEntertainmentTrending  = (ProductTypeBookEntertainment * 100) + iota + 1 //热门1401
	ProductTypeBookEntertainmentnovels                                                      //小说1402
	ProductTypeBookEntertainmentcomics                                                      //漫画1403
	ProductTypeBookEntertainmentmagazines                                                   //杂志1404
	ProductTypeBookEntertainmentmusic                                                       //音乐1405
	ProductTypeBookEntertainmentmovies                                                      //电影1406
)

const (
	//二级分类（钟表眼镜）
	ProductTypeClockGlassesTrending = (ProductTypeClockGlasses * 100) + iota + 1 //热门1501
	ProductTypeClockGlasseswatches                                               //表带1502
	ProductTypeClockGlassesclocks                                                //手表1503
	ProductTypeClockGlassesglasses                                               //眼镜1504
)

const (
	//二级分类（工业品）
	ProductTypeIndustrialproductsTrending     = (ProductTypeIndustrialproducts * 100) + iota + 1 //热门1601
	ProductTypeIndustrialproductstools                                                           //工具1602
	ProductTypeIndustrialproductscables                                                          //清洁用品1603
	ProductTypeIndustrialproductscontrolpanel                                                    //控制器1604
)

const (
	//二级分类（母婴童装）
	ProductTypeMkwTrending     = (ProductTypeMkw * 100) + iota + 1 //热门1701
	ProductTypeMkwbabyproducts                                     //婴童用品1702
	ProductTypeMkwtoys                                             //玩具1703
	ProductTypeMkwbabydiapers                                      //奶粉1704
	ProductTypeMkwclothing                                         //服装1705
	ProductTypeMkwbabyshoes                                        //童鞋1706
)

const (
	//二级分类（黄金珠宝）
	ProductTypeGoldJewelryTrending  = (ProductTypeGoldJewelry * 100) + iota + 1 //热门1801
	ProductTypeGoldJewelrygold                                                  //黄金1802
	ProductTypeGoldJewelryjewelry                                               //珠宝1803
	ProductTypeGoldJewelryrings                                                 //戒指1804
	ProductTypeGoldJewelrynecklaces                                             //项链1805
)

const (
	//二级分类（家具家装）
	ProductTypeFurnitureTrending = (ProductTypeFurniture * 100) + iota + 1 //热门1901
	ProductTypeFurnitureschair                                             //客厅1902
	ProductTypeFurnitureshelves                                            //书架1903
	ProductTypeFurnituretable                                              //柜子1904
)

const (
	//二级分类（厨具）
	ProductTypeKitchenwareTrending          = (ProductTypeKitchenware * 100) + iota + 1 //热门2001
	ProductTypeKitchenwarecutlery                                                       //刀剪菜板2002
	ProductTypeKitchenwarecookware                                                      //厨房用具2003
	ProductTypeKitchenwareutensils                                                      //餐具2004
	ProductTypeKitchenwarecookingutensils                                               //锅2005
	ProductTypeKitchenwarekitchenappliances                                             //厨房电器2006
)

const (
	//二级分类（元器件）
	ProductTypeComponentsTrending     = (ProductTypeComponents * 100) + iota + 1 //热门2101
	ProductTypeComponentscables                                                  //电缆2102
	ProductTypeComponentssensors                                                 //传感器2103
	ProductTypeComponentscontrolpanel                                            //控制器2104
	ProductTypeComponentscontrolunit                                             //控制单元2105
	ProductTypeComponentscontrolboard                                            //控制板2106
)

const (
	//二级分类（文具）
	ProductTypeStationeryTrending = (ProductTypeStationery * 100) + iota + 1 //热门2201
	ProductTypeStationerypen                                                 //笔2202
	ProductTypeStationerypencil                                              //铅笔2203
	ProductTypeStationeryruler                                               //尺子2204
	ProductTypeStationeryeraser                                              //橡皮2205
	ProductTypeStationerypaper                                               //纸2206
	ProductTypeStationeryink                                                 //墨水2207
)

const (
	//二级分类（箱包）
	ProductTypeLuggageTrending  = (ProductTypeLuggage * 100) + iota + 1 //热门2301
	ProductTypeLuggagebags                                              //包包2302
	ProductTypeLuggagebackpacks                                         //手提包2303
	ProductTypeLuggagewallets                                           //钱包2304
	ProductTypeLuggagesuitcases                                         //行李箱2305
)

const (
	//二级分类（男装）
	ProductTypeMenswearTrending = (ProductTypeMenswear * 100) + iota + 1 //热门2401
	ProductTypeMenswearclothing                                          //男装服装2402
	ProductTypeMenswearshoes                                             //男装鞋2403
	ProductTypeMenswearjewelry                                           //男装饰品2404
	ProductTypeMenswearbags                                              //男包2405
)

const (
	//二级分类（美容护肤）
	ProductTypeBeautyskinTrending = (ProductTypeBeautyskin * 100) + iota + 1 //热门2501
	ProductTypeBeautyskinface                                                //面部护肤2502
	ProductTypeBeautyskinbody                                                //身体护肤2503
	ProductTypeBeautyskinmakeup                                              //美妆2504
	ProductTypeBeautyskinskincare                                            //皮肤护理2505
)

const (
	//二级分类（家居）
	ProductTypeHouseholdTrending   = (ProductTypeHousehold * 100) + iota + 1 //热门2601
	ProductTypeHouseholdfurniture                                            //家具2602
	ProductTypeHouseholddecorative                                           //装饰2603
	ProductTypeHouseholdlighting                                             //照明2604
	ProductTypeHouseholdappliances                                           //家电2605
)

const (
	//二级分类（电脑办公）
	ProductTypeOfficeworkcomputersTrending    = (ProductTypeOfficeworkcomputers * 100) + iota + 1 //热门2701
	ProductTypeOfficeworkcomputersdesktop                                                         //台式机2702
	ProductTypeOfficeworkcomputerslaptop                                                          //笔记本2703
	ProductTypeOfficeworkcomputersprinter                                                         //打印机2704
	ProductTypeOfficeworkcomputersaccessories                                                     //配件2705
	ProductTypeOfficeworkcomputerssoftware                                                        //软件2706
	ProductTypeOfficeworkcomputersnetwork                                                         //网络设备2707
)

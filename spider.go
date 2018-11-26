package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var firAreaName = []string{"臺北市", "新北市", "桃園市", "臺中市", "臺南市", "高雄市", "新竹縣", "苗栗縣", "彰化縣", "南投縣", "雲林縣", "嘉義縣", "屏東縣", "宜蘭縣", "花蓮縣", "臺東縣", "澎湖縣", "金門縣", "連江縣", "基隆市", "新竹市", "嘉義市"}
var firAreaID = []string{"63000", "65000", "68000", "66000", "67000", "64000", "10004", "10005", "10007", "10008", "10009", "10010", "10013", "10002", "10015", "10014", "10016", "09020", "09007", "10017", "10018", "10020"}
var thrAreaNameJsonString = "[[\"松山區\",\"信義區\",\"大安區\",\"中山區\",\"中正區\",\"大同區\",\"萬華區\",\"文山區\",\"南港區\",\"內湖區\",\"士林區\",\"北投區\"],[\"板橋區\",\"三重區\",\"中和區\",\"永和區\",\"新莊區\",\"新店區\",\"樹林區\",\"鶯歌區\",\"三峽區\",\"淡水區\",\"汐止區\",\"瑞芳區\",\"土城區\",\"蘆洲區\",\"五股區\",\"泰山區\",\"林口區\",\"深坑區\",\"石碇區\",\"坪林區\",\"三芝區\",\"石門區\",\"八里區\",\"平溪區\",\"雙溪區\",\"貢寮區\",\"金山區\",\"萬里區\",\"烏來區\"],[\"桃園區\",\"中壢區\",\"大溪區\",\"楊梅區\",\"蘆竹區\",\"大園區\",\"龜山區\",\"八德區\",\"龍潭區\",\"平鎮區\",\"新屋區\",\"觀音區\",\"復興區\"],[\"中區\",\"東區\",\"南區\",\"西區\",\"北區\",\"西屯區\",\"南屯區\",\"北屯區\",\"豐原區\",\"東勢區\",\"大甲區\",\"清水區\",\"沙鹿區\",\"梧棲區\",\"后里區\",\"神岡區\",\"潭子區\",\"大雅區\",\"新社區\",\"石岡區\",\"外埔區\",\"大安區\",\"烏日區\",\"大肚區\",\"龍井區\",\"霧峰區\",\"太平區\",\"大里區\",\"和平區\"],[\"新營區\",\"鹽水區\",\"白河區\",\"柳營區\",\"後壁區\",\"東山區\",\"麻豆區\",\"下營區\",\"六甲區\",\"官田區\",\"大內區\",\"佳里區\",\"學甲區\",\"西港區\",\"七股區\",\"將軍區\",\"北門區\",\"新化區\",\"善化區\",\"新市區\",\"安定區\",\"山上區\",\"玉井區\",\"楠西區\",\"南化區\",\"左鎮區\",\"仁德區\",\"歸仁區\",\"關廟區\",\"龍崎區\",\"永康區\",\"東區\",\"南區\",\"北區\",\"安南區\",\"安平區\",\"中西區\"],[\"鹽埕區\",\"鼓山區\",\"左營區\",\"楠梓區\",\"三民區\",\"新興區\",\"前金區\",\"苓雅區\",\"前鎮區\",\"旗津區\",\"小港區\",\"鳳山區\",\"林園區\",\"大寮區\",\"大樹區\",\"大社區\",\"仁武區\",\"鳥松區\",\"岡山區\",\"橋頭區\",\"燕巢區\",\"田寮區\",\"阿蓮區\",\"路竹區\",\"湖內區\",\"茄萣區\",\"永安區\",\"彌陀區\",\"梓官區\",\"旗山區\",\"美濃區\",\"六龜區\",\"甲仙區\",\"杉林區\",\"內門區\",\"茂林區\",\"桃源區\",\"那瑪夏區\"],[\"竹北市\",\"竹東鎮\",\"新埔鎮\",\"關西鎮\",\"湖口鄉\",\"新豐鄉\",\"芎林鄉\",\"橫山鄉\",\"北埔鄉\",\"寶山鄉\",\"峨眉鄉\",\"尖石鄉\",\"五峰鄉\"],[\"苗栗市\",\"苑裡鎮\",\"通霄鎮\",\"竹南鎮\",\"頭份市\",\"後龍鎮\",\"卓蘭鎮\",\"大湖鄉\",\"公館鄉\",\"銅鑼鄉\",\"南庄鄉\",\"頭屋鄉\",\"三義鄉\",\"西湖鄉\",\"造橋鄉\",\"三灣鄉\",\"獅潭鄉\",\"泰安鄉\"],[\"彰化市\",\"鹿港鎮\",\"和美鎮\",\"線西鄉\",\"伸港鄉\",\"福興鄉\",\"秀水鄉\",\"花壇鄉\",\"芬園鄉\",\"員林市\",\"溪湖鎮\",\"田中鎮\",\"大村鄉\",\"埔鹽鄉\",\"埔心鄉\",\"永靖鄉\",\"社頭鄉\",\"二水鄉\",\"北斗鎮\",\"二林鎮\",\"田尾鄉\",\"埤頭鄉\",\"芳苑鄉\",\"大城鄉\",\"竹塘鄉\",\"溪州鄉\"],[\"南投市\",\"埔里鎮\",\"草屯鎮\",\"竹山鎮\",\"集集鎮\",\"名間鄉\",\"鹿谷鄉\",\"中寮鄉\",\"魚池鄉\",\"國姓鄉\",\"水里鄉\",\"信義鄉\",\"仁愛鄉\"],[\"斗六市\",\"斗南鎮\",\"虎尾鎮\",\"西螺鎮\",\"土庫鎮\",\"北港鎮\",\"古坑鄉\",\"大埤鄉\",\"莿桐鄉\",\"林內鄉\",\"二崙鄉\",\"崙背鄉\",\"麥寮鄉\",\"東勢鄉\",\"褒忠鄉\",\"臺西鄉\",\"元長鄉\",\"四湖鄉\",\"口湖鄉\",\"水林鄉\"],[\"太保市\",\"朴子市\",\"布袋鎮\",\"大林鎮\",\"民雄鄉\",\"溪口鄉\",\"新港鄉\",\"六腳鄉\",\"東石鄉\",\"義竹鄉\",\"鹿草鄉\",\"水上鄉\",\"中埔鄉\",\"竹崎鄉\",\"梅山鄉\",\"番路鄉\",\"大埔鄉\",\"阿里山鄉\"],[\"屏東市\",\"潮州鎮\",\"東港鎮\",\"恆春鎮\",\"萬丹鄉\",\"長治鄉\",\"麟洛鄉\",\"九如鄉\",\"里港鄉\",\"鹽埔鄉\",\"高樹鄉\",\"萬巒鄉\",\"內埔鄉\",\"竹田鄉\",\"新埤鄉\",\"枋寮鄉\",\"新園鄉\",\"崁頂鄉\",\"林邊鄉\",\"南州鄉\",\"佳冬鄉\",\"琉球鄉\",\"車城鄉\",\"滿州鄉\",\"枋山鄉\",\"三地門鄉\",\"霧臺鄉\",\"瑪家鄉\",\"泰武鄉\",\"來義鄉\",\"春日鄉\",\"獅子鄉\",\"牡丹鄉\"],[\"宜蘭市\",\"羅東鎮\",\"蘇澳鎮\",\"頭城鎮\",\"礁溪鄉\",\"壯圍鄉\",\"員山鄉\",\"冬山鄉\",\"五結鄉\",\"三星鄉\",\"大同鄉\",\"南澳鄉\"],[\"花蓮市\",\"鳳林鎮\",\"玉里鎮\",\"新城鄉\",\"吉安鄉\",\"壽豐鄉\",\"光復鄉\",\"豐濱鄉\",\"瑞穗鄉\",\"富里鄉\",\"秀林鄉\",\"萬榮鄉\",\"卓溪鄉\"],[\"臺東市\",\"成功鎮\",\"關山鎮\",\"卑南鄉\",\"鹿野鄉\",\"池上鄉\",\"東河鄉\",\"長濱鄉\",\"太麻里鄉\",\"大武鄉\",\"綠島鄉\",\"海端鄉\",\"延平鄉\",\"金峰鄉\",\"達仁鄉\",\"蘭嶼鄉\"],[\"馬公市\",\"湖西鄉\",\"白沙鄉\",\"西嶼鄉\",\"望安鄉\",\"七美鄉\"],[\"金城鎮\",\"金沙鎮\",\"金湖鎮\",\"金寧鄉\",\"烈嶼鄉\",\"烏坵鄉\"],[\"南竿鄉\",\"北竿鄉\",\"莒光鄉\",\"東引鄉\"],[\"中正區\",\"七堵區\",\"暖暖區\",\"仁愛區\",\"中山區\",\"安樂區\",\"信義區\"],[\"東區\",\"北區\",\"香山區\"],[\"東區\",\"西區\"]]"
var thrAreaIDJsonString = "[[\"6300000010\",\"6300000020\",\"6300000030\",\"6300000040\",\"6300000050\",\"6300000060\",\"6300000070\",\"6300000080\",\"6300000090\",\"6300000100\",\"6300000110\",\"6300000120\"],[\"6500000010\",\"6500000020\",\"6500000030\",\"6500000040\",\"6500000050\",\"6500000060\",\"6500000070\",\"6500000080\",\"6500000090\",\"6500000100\",\"6500000110\",\"6500000120\",\"6500000130\",\"6500000140\",\"6500000150\",\"6500000160\",\"6500000170\",\"6500000180\",\"6500000190\",\"6500000200\",\"6500000210\",\"6500000220\",\"6500000230\",\"6500000240\",\"6500000250\",\"6500000260\",\"6500000270\",\"6500000280\",\"6500000290\"],[\"6800000010\",\"6800000020\",\"6800000030\",\"6800000040\",\"6800000050\",\"6800000060\",\"6800000070\",\"6800000080\",\"6800000090\",\"6800000100\",\"6800000110\",\"6800000120\",\"6800000130\"],[\"6600000010\",\"6600000020\",\"6600000030\",\"6600000040\",\"6600000050\",\"6600000060\",\"6600000070\",\"6600000080\",\"6600000090\",\"6600000100\",\"6600000110\",\"6600000120\",\"6600000130\",\"6600000140\",\"6600000150\",\"6600000160\",\"6600000170\",\"6600000180\",\"6600000190\",\"6600000200\",\"6600000210\",\"6600000220\",\"6600000230\",\"6600000240\",\"6600000250\",\"6600000260\",\"6600000270\",\"6600000280\",\"6600000290\"],[\"6700000010\",\"6700000020\",\"6700000030\",\"6700000040\",\"6700000050\",\"6700000060\",\"6700000070\",\"6700000080\",\"6700000090\",\"6700000100\",\"6700000110\",\"6700000120\",\"6700000130\",\"6700000140\",\"6700000150\",\"6700000160\",\"6700000170\",\"6700000180\",\"6700000190\",\"6700000200\",\"6700000210\",\"6700000220\",\"6700000230\",\"6700000240\",\"6700000250\",\"6700000260\",\"6700000270\",\"6700000280\",\"6700000290\",\"6700000300\",\"6700000310\",\"6700000320\",\"6700000330\",\"6700000340\",\"6700000350\",\"6700000360\",\"6700000370\"],[\"6400000010\",\"6400000020\",\"6400000030\",\"6400000040\",\"6400000050\",\"6400000060\",\"6400000070\",\"6400000080\",\"6400000090\",\"6400000100\",\"6400000110\",\"6400000120\",\"6400000130\",\"6400000140\",\"6400000150\",\"6400000160\",\"6400000170\",\"6400000180\",\"6400000190\",\"6400000200\",\"6400000210\",\"6400000220\",\"6400000230\",\"6400000240\",\"6400000250\",\"6400000260\",\"6400000270\",\"6400000280\",\"6400000290\",\"6400000300\",\"6400000310\",\"6400000320\",\"6400000330\",\"6400000340\",\"6400000350\",\"6400000360\",\"6400000370\",\"6400000380\"],[\"1000400010\",\"1000400020\",\"1000400030\",\"1000400040\",\"1000400050\",\"1000400060\",\"1000400070\",\"1000400080\",\"1000400090\",\"1000400100\",\"1000400110\",\"1000400120\",\"1000400130\"],[\"1000500010\",\"1000500020\",\"1000500030\",\"1000500040\",\"1000500050\",\"1000500060\",\"1000500070\",\"1000500080\",\"1000500090\",\"1000500100\",\"1000500110\",\"1000500120\",\"1000500130\",\"1000500140\",\"1000500150\",\"1000500160\",\"1000500170\",\"1000500180\"],[\"1000700010\",\"1000700020\",\"1000700030\",\"1000700040\",\"1000700050\",\"1000700060\",\"1000700070\",\"1000700080\",\"1000700090\",\"1000700100\",\"1000700110\",\"1000700120\",\"1000700130\",\"1000700140\",\"1000700150\",\"1000700160\",\"1000700170\",\"1000700180\",\"1000700190\",\"1000700200\",\"1000700210\",\"1000700220\",\"1000700230\",\"1000700240\",\"1000700250\",\"1000700260\"],[\"1000800010\",\"1000800020\",\"1000800030\",\"1000800040\",\"1000800050\",\"1000800060\",\"1000800070\",\"1000800080\",\"1000800090\",\"1000800100\",\"1000800110\",\"1000800120\",\"1000800130\"],[\"1000900010\",\"1000900020\",\"1000900030\",\"1000900040\",\"1000900050\",\"1000900060\",\"1000900070\",\"1000900080\",\"1000900090\",\"1000900100\",\"1000900110\",\"1000900120\",\"1000900130\",\"1000900140\",\"1000900150\",\"1000900160\",\"1000900170\",\"1000900180\",\"1000900190\",\"1000900200\"],[\"1001000010\",\"1001000020\",\"1001000030\",\"1001000040\",\"1001000050\",\"1001000060\",\"1001000070\",\"1001000080\",\"1001000090\",\"1001000100\",\"1001000110\",\"1001000120\",\"1001000130\",\"1001000140\",\"1001000150\",\"1001000160\",\"1001000170\",\"1001000180\"],[\"1001300010\",\"1001300020\",\"1001300030\",\"1001300040\",\"1001300050\",\"1001300060\",\"1001300070\",\"1001300080\",\"1001300090\",\"1001300100\",\"1001300110\",\"1001300120\",\"1001300130\",\"1001300140\",\"1001300150\",\"1001300160\",\"1001300170\",\"1001300180\",\"1001300190\",\"1001300200\",\"1001300210\",\"1001300220\",\"1001300230\",\"1001300240\",\"1001300250\",\"1001300260\",\"1001300270\",\"1001300280\",\"1001300290\",\"1001300300\",\"1001300310\",\"1001300320\",\"1001300330\"],[\"1000200010\",\"1000200020\",\"1000200030\",\"1000200040\",\"1000200050\",\"1000200060\",\"1000200070\",\"1000200080\",\"1000200090\",\"1000200100\",\"1000200110\",\"1000200120\"],[\"1001500010\",\"1001500020\",\"1001500030\",\"1001500040\",\"1001500050\",\"1001500060\",\"1001500070\",\"1001500080\",\"1001500090\",\"1001500100\",\"1001500110\",\"1001500120\",\"1001500130\"],[\"1001400010\",\"1001400020\",\"1001400030\",\"1001400040\",\"1001400050\",\"1001400060\",\"1001400070\",\"1001400080\",\"1001400090\",\"1001400100\",\"1001400110\",\"1001400120\",\"1001400130\",\"1001400140\",\"1001400150\",\"1001400160\"],[\"1001600010\",\"1001600020\",\"1001600030\",\"1001600040\",\"1001600050\",\"1001600060\"],[\"0902000010\",\"0902000020\",\"0902000030\",\"0902000040\",\"0902000050\",\"0902000060\"],[\"0900700010\",\"0900700020\",\"0900700030\",\"0900700040\"],[\"1001700010\",\"1001700020\",\"1001700030\",\"1001700040\",\"1001700050\",\"1001700060\",\"1001700070\"],[\"1001800010\",\"1001800020\",\"1001800030\"],[\"1002000010\",\"1002000020\"]]"
var thrAreaName [][]string = [][]string{}
var thrAreaID [][]string = [][]string{}
var thrAreaStationString = "[[[\"527,643\"],[\"644,785\"],[\"1214,1401\"],[\"786,918\"],[\"1001,1096\"],[\"919,1000\"],[\"1097,1213\"],[\"1402,1563\"],[\"455,526\"],[\"305,454\"],[\"142,304\"],[\"1,141\"]],[[\"894,1231\"],[\"659,893\"],[\"1232,1475\"],[\"1476,1613\"],[\"302,554\"],[\"1987,2162\"],[\"1614,1724\"],[\"1725,1774\"],[\"1917,1986\"],[\"29,124\"],[\"2336,2446\"],[\"2220,2265\"],[\"1775,1916\"],[\"555,658\"],[\"207,254\"],[\"255,301\"],[\"145,206\"],[\"2163,2178\"],[\"2179,2191\"],[\"2192,2201\"],[\"11,28\"],[\"1,10\"],[\"125,144\"],[\"2208,2219\"],[\"2266,2282\"],[\"2283,2300\"],[\"2301,2319\"],[\"2320,2335\"],[\"2202,2207\"]],[[\"1,216\"],[\"591,800\"],[\"522,580\"],[\"926,1007\"],[\"396,471\"],[\"472,521\"],[\"217,294\"],[\"295,395\"],[\"1008,1071\"],[\"801,925\"],[\"1072,1101\"],[\"1102,1142\"],[\"581,590\"]],[[\"1070,1081\"],[\"1148,1196\"],[\"1197,1270\"],[\"1082,1147\"],[\"983,1069\"],[\"625,744\"],[\"745,845\"],[\"846,982\"],[\"394,481\"],[\"1513,1551\"],[\"1,54\"],[\"91,149\"],[\"186,238\"],[\"150,185\"],[\"357,393\"],[\"482,523\"],[\"570,624\"],[\"524,569\"],[\"1562,1584\"],[\"1552,1561\"],[\"71,90\"],[\"55,70\"],[\"317,356\"],[\"284,316\"],[\"239,283\"],[\"1470,1512\"],[\"1271,1363\"],[\"1364,1469\"],[\"1585,1601\"]],[[\"116,164\"],[\"85,115\"],[\"30,61\"],[\"165,189\"],[\"1,29\"],[\"62,84\"],[\"381,415\"],[\"343,361\"],[\"362,380\"],[\"416,435\"],[\"436,446\"],[\"282,322\"],[\"208,233\"],[\"323,342\"],[\"258,281\"],[\"234,257\"],[\"190,207\"],[\"581,617\"],[\"492,520\"],[\"545,571\"],[\"521,544\"],[\"572,580\"],[\"466,481\"],[\"447,455\"],[\"456,465\"],[\"482,491\"],[\"1188,1236\"],[\"1237,1282\"],[\"1283,1314\"],[\"1315,1322\"],[\"708,849\"],[\"1095,1187\"],[\"1015,1094\"],[\"850,930\"],[\"618,707\"],[\"983,1014\"],[\"931,982\"]],[[\"845,868\"],[\"768,844\"],[\"534,636\"],[\"435,533\"],[\"887,1075\"],[\"1098,1139\"],[\"1076,1097\"],[\"1140,1255\"],[\"1484,1602\"],[\"869,886\"],[\"1603,1691\"],[\"1256,1483\"],[\"1773,1823\"],[\"1692,1772\"],[\"736,767\"],[\"637,660\"],[\"661,707\"],[\"708,735\"],[\"261,332\"],[\"404,434\"],[\"333,357\"],[\"237,248\"],[\"216,236\"],[\"175,215\"],[\"151,174\"],[\"130,150\"],[\"249,260\"],[\"358,375\"],[\"376,403\"],[\"60,94\"],[\"95,126\"],[\"20,31\"],[\"13,19\"],[\"32,41\"],[\"42,59\"],[\"127,129\"],[\"1,8\"],[\"9,12\"]],[[\"1,103\"],[\"104,175\"],[\"176,205\"],[\"206,232\"],[\"233,279\"],[\"280,318\"],[\"319,337\"],[\"338,353\"],[\"354,363\"],[\"364,375\"],[\"376,383\"],[\"384,398\"],[\"399,404\"]],[[\"1,67\"],[\"183,222\"],[\"149,182\"],[\"270,330\"],[\"331,401\"],[\"223,257\"],[\"442,460\"],[\"420,434\"],[\"68,95\"],[\"108,123\"],[\"411,419\"],[\"96,107\"],[\"124,139\"],[\"140,148\"],[\"258,269\"],[\"402,410\"],[\"435,441\"],[\"461,470\"]],[[\"1,184\"],[\"243,307\"],[\"373,440\"],[\"469,480\"],[\"441,468\"],[\"308,340\"],[\"341,372\"],[\"185,221\"],[\"222,242\"],[\"481,578\"],[\"640,681\"],[\"742,777\"],[\"579,605\"],[\"682,710\"],[\"711,741\"],[\"606,639\"],[\"778,815\"],[\"816,832\"],[\"833,862\"],[\"942,980\"],[\"863,887\"],[\"888,914\"],[\"998,1033\"],[\"981,997\"],[\"1034,1049\"],[\"915,941\"]],[[\"1,84\"],[\"85,152\"],[\"153,228\"],[\"229,281\"],[\"282,294\"],[\"295,331\"],[\"332,352\"],[\"353,376\"],[\"377,401\"],[\"402,421\"],[\"422,443\"],[\"444,463\"],[\"464,488\"]],[[\"411,479\"],[\"498,530\"],[\"109,154\"],[\"340,373\"],[\"80,108\"],[\"256,296\"],[\"531,557\"],[\"480,497\"],[\"374,395\"],[\"396,410\"],[\"317,339\"],[\"297,316\"],[\"1,28\"],[\"51,65\"],[\"66,79\"],[\"29,50\"],[\"180,204\"],[\"155,179\"],[\"205,229\"],[\"230,255\"]],[[\"1,33\"],[\"233,267\"],[\"325,353\"],[\"169,193\"],[\"93,137\"],[\"194,210\"],[\"138,168\"],[\"268,297\"],[\"298,324\"],[\"354,379\"],[\"34,53\"],[\"54,92\"],[\"380,413\"],[\"414,449\"],[\"211,232\"],[\"450,464\"],[\"465,473\"],[\"474,486\"]],[[\"1,149\"],[\"150,189\"],[\"190,227\"],[\"228,252\"],[\"253,291\"],[\"292,313\"],[\"314,323\"],[\"324,337\"],[\"338,357\"],[\"358,377\"],[\"378,399\"],[\"400,417\"],[\"418,457\"],[\"458,473\"],[\"474,485\"],[\"486,507\"],[\"508,535\"],[\"536,547\"],[\"548,563\"],[\"564,574\"],[\"575,590\"],[\"591,599\"],[\"600,612\"],[\"613,620\"],[\"621,626\"],[\"627,636\"],[\"637,644\"],[\"645,652\"],[\"653,659\"],[\"660,668\"],[\"669,674\"],[\"675,682\"],[\"683,689\"]],[[\"1,78\"],[\"79,128\"],[\"129,171\"],[\"172,203\"],[\"204,235\"],[\"236,258\"],[\"259,289\"],[\"290,328\"],[\"329,361\"],[\"362,382\"],[\"383,392\"],[\"393,399\"]],[[\"1,81\"],[\"195,209\"],[\"254,276\"],[\"82,100\"],[\"115,173\"],[\"174,194\"],[\"210,226\"],[\"227,233\"],[\"241,253\"],[\"277,289\"],[\"101,114\"],[\"234,240\"],[\"290,299\"]],[[\"1,69\"],[\"70,86\"],[\"87,95\"],[\"96,117\"],[\"118,127\"],[\"128,139\"],[\"140,150\"],[\"151,161\"],[\"162,176\"],[\"177,182\"],[\"183,186\"],[\"187,192\"],[\"193,197\"],[\"198,202\"],[\"203,209\"],[\"210,213\"]],[[\"1,50\"],[\"51,73\"],[\"74,88\"],[\"89,100\"],[\"101,110\"],[\"111,116\"]],[[\"1,21\"],[\"66,78\"],[\"49,65\"],[\"22,38\"],[\"39,46\"],[\"47,48\"]],[[\"1,4\"],[\"5,6\"],[\"7,8\"],[\"9,9\"]],[[\"1,35\"],[\"219,257\"],[\"191,218\"],[\"72,106\"],[\"107,141\"],[\"142,190\"],[\"36,71\"]],[[\"1,139\"],[\"140,244\"],[\"245,299\"]],[[\"1,81\"],[\"82,177\"]]]"
var thrAreaStation [][][]string = [][][]string{}

var thrAreaUrl [][][]string = [][][]string{}
var thrAreaStationNumber [][][]string = [][][]string{}

var title = []string{
	"縣市",
	"鄉鎮市區",
	"開票所",
	"同意票",
	"不同意票",
	"有效票",
	"無效票",
	"投票數",
	"投票權人數",
	"投票率",
	"投票權人數百分比",
	"總數",
}

var total = 15887
var nowCount = 0

func main() {
	parseJsons()
	makeAreaUrl()
	// fmt.Println(thrAreaUrl

	// getResult(6, thrAreaUrl[0][0][0])

	// result := getUrl("http://referendum.2018.nat.gov.tw/pc/zh_TW/08/67000000100000000.html")
	// getTrTFromHTML(result)
	topic := 9
	f, err := os.OpenFile(fmt.Sprintf("%d.csv", topic), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	nowCount = 0
	// runResult(topic, 3, f) // for testing
	runResult(topic, 10000, f) // for testing

	f.Close()
}

func runResult(topic int, limit int, out io.Writer) {

	w := csv.NewWriter(out)
	w.Write(title)
	for i, list := range thrAreaUrl {
		if i >= limit {
			break
		}
		for j, item := range list {
			if j >= limit {
				break
			}

			for k, _ := range item {
				if k >= limit {
					break
				}

				// for j, _ := range list {
				// for k, _ := range
				result := getResult(topic, thrAreaUrl[i][j][k])
				label := []string{
					firAreaName[i],
					thrAreaName[i][j],
					"",
				}
				// fmt.Println(result)
				label[2] = thrAreaStationNumber[i][j][k]
				w.Write(append(label, result...))
				nowCount = nowCount + 1
				fmt.Println("nowCount:", nowCount, total)
				time.Sleep(30 * time.Millisecond)
			}

		}
	}
	w.Flush()

}

func parseJsons() {
	json.Unmarshal([]byte(thrAreaNameJsonString), &thrAreaName)
	json.Unmarshal([]byte(thrAreaIDJsonString), &thrAreaID)
	json.Unmarshal([]byte(thrAreaStationString), &thrAreaStation)
}

func makeAreaUrl() {
	counter := 0
	for i, list := range thrAreaStation {
		cAreaUrl := [][]string{}
		cthrAreaStationNumber := [][]string{}
		for j, item := range list {

			dAreaUrl := []string{}
			dthrAreaStationNumber := []string{}
			// fmt.Println(i, j, item[0])
			ar := strings.Split(item[0], ",")
			start, _ := strconv.ParseInt(ar[0], 10, 16)
			end, _ := strconv.ParseInt(ar[1], 10, 16)
			// fmt.Println("s", start, end)
			// 63000000200000644.html
			areaId := thrAreaID[i][j]
			for k := start; k <= end; k++ {
				str := fmt.Sprintf("%s000%04d.html", areaId, k)
				dAreaUrl = append(dAreaUrl, str)
				dthrAreaStationNumber = append(dthrAreaStationNumber, fmt.Sprintf("%04d", k))
				counter = counter + 1
			}
			cAreaUrl = append(cAreaUrl, dAreaUrl)
			cthrAreaStationNumber = append(cthrAreaStationNumber, dthrAreaStationNumber)

		}

		thrAreaUrl = append(thrAreaUrl, cAreaUrl)
		thrAreaStationNumber = append(thrAreaStationNumber, cthrAreaStationNumber)
	}
	fmt.Println("total:", counter)
}

func getResult(topicId int, url string) []string {
	page := getUrl(fmt.Sprintf("http://referendum.2018.nat.gov.tw/pc/zh_TW/%02d/%s", topicId, url))
	同意票數, 不同意票數, 有效票數, 無效票數, 投票數, 投票權人數, 投票率, 投票權人數百分比, 總數 := getResultTFromHTML(page)
	// fmt.Println(同意票數, 不同意票數, 有效票數, 無效票數, 投票數, 投票權人數, 投票率, 投票權人數百分比, 總數)
	return []string{同意票數, 不同意票數, 有效票數, 無效票數, 投票數, 投票權人數, 投票率, 投票權人數百分比, 總數}
	// fmt.Println(page)

}

func getUrl(path string) string {

	req, _ := http.NewRequest(http.MethodGet, path, nil)

	resp, _ := http.DefaultClient.Do(req)

	data, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(string(data))
	return string(data)
	// strings.

}

func getResultTFromHTML(input string) (同意票數 string, 不同意票數 string, 有效票數 string, 無效票數 string, 投票數 string, 投票權人數 string, 投票率 string, 投票權人數百分比 string, 總數 string) {
	// text := strings.Replace(html, "\n", " ", -1)

	// reg, _ := regexp.Compile("<tr class=\"trT\"> *<td class(.*)? *</tr>")
	// res := reg.FindAllStringSubmatch(text, -1)
	// fmt.Println(res)
	z := html.NewTokenizer(strings.NewReader(input))
	fmt.Println("zz", z)
	第一欄標題 := "未知"
	for {
		tt := z.Next()
		// fmt.Println("tt", string(name))
		if tt == html.ErrorToken {
			// ...
			return
		} else if tt == html.StartTagToken {

			name, hasAttr := z.TagName()
			// fmt.Println("tt", string(name))
			if string(name) == "tr" {
				fmt.Println("text: ", string(z.Raw()))

				if hasAttr {
					k, a, ma := z.TagAttr()
					fmt.Println("tag:", string(k), string(a), ma)
					if string(a) == "trT" {
						// start
						if 第一欄標題 == "投票數" {
							投票數, 投票權人數, 投票率, 投票權人數百分比 = parse投票數表(z)
						} else if 第一欄標題 == "同意票數" {
							同意票數, 不同意票數, 有效票數, 無效票數 = parse同意票數(z)
						}

					} else if string(a) == "trHeaderT" {
						// set table 1 or 2

						// fmt.Println("trHeaderT tag:", string(k), string(a), ma)
						z.Next()
						z.Next()
						z.Next()
						第一欄標題 = string(z.Raw())
						fmt.Println("第一欄標題:", 第一欄標題)
						// if 第一欄標題 == "同意票數" {

						// }else if
					} else if string(a) == "trFooterT" {
						z.Next()
						z.Next()
						z.Next()
						總數s := string(z.Raw())
						fmt.Println("總數:", 總數s)
						reg := regexp.MustCompile("(\\d+)/(\\d+)")
						x := reg.FindStringSubmatch(總數s)
						總數 = x[1] + ""
						fmt.Println("投開票所數量:", x[1])
					}

				}

			}

		}
		// Process the current token.

	}
}

func parseForFourColumnRow(z *html.Tokenizer) (string, string, string, string) {
	z.Next()

	// fetch td
	z.Next()
	z.Next()
	投票數 := string(z.Raw()) // number of vote
	投票數 = strings.Replace(投票數, ",", "", -1)
	// fmt.Println("投票數:", 投票數)
	z.Next()
	z.Next()

	z.Next()
	z.Next()
	投票權人數 := string(z.Raw())
	投票權人數 = strings.Replace(投票權人數, ",", "", -1)
	// fmt.Println("投票權人數:", 投票權人數)
	z.Next()
	z.Next()

	z.Next()
	z.Next()
	投票率 := string(z.Raw()) // number of vote
	投票率 = strings.Replace(投票率, ",", "", -1)
	// fmt.Println("投票率:", 投票率)
	z.Next()
	z.Next()

	z.Next()
	z.Next()
	投票權人數百分比 := string(z.Raw()) // number of vote
	投票權人數百分比 = strings.Replace(投票權人數百分比, ",", "", -1)
	// fmt.Println("投票權人數百分比:", 投票權人數百分比)
	z.Next()
	z.Next()
	return 投票數, 投票權人數, 投票率, 投票權人數百分比

}

func parse同意票數(z *html.Tokenizer) (string, string, string, string) {
	同意票數, 不同意票數, 有效票數, 無效票數 := parseForFourColumnRow(z)
	fmt.Println("同意票數:", 同意票數)
	fmt.Println("不同意票數:", 不同意票數)
	fmt.Println("有效票數:", 有效票數)
	fmt.Println("無效票數:", 無效票數)

	return 同意票數, 不同意票數, 有效票數, 無效票數
}

func parse投票數表(z *html.Tokenizer) (string, string, string, string) {
	投票數, 投票權人數, 投票率, 投票權人數百分比 := parseForFourColumnRow(z)
	fmt.Println("投票數:", 投票數)
	fmt.Println("投票權人數:", 投票權人數)
	fmt.Println("投票率:", 投票率)
	fmt.Println("投票權人數百分比:", 投票權人數百分比)

	return 投票數, 投票權人數, 投票率, 投票權人數百分比

}

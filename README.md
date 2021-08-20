# stockinfo

## 簡介:
獲取台股每日收盤資料(個股相關交易資訊,三大法人交易資訊等)  
並根據輸入條件篩選符合條件的個股


## 使用環境
DBMS:mysql  
web framework: beego  
orm framework: gorm    

## 使用前置
1.開啟mysql 
2.建立資料庫  
3.開啟app.conf進行DB設定  
4.運行項目(自動建置table)   
5.運行sql目錄內的腳本,將資料導入DB內  

## 操作說明
啟動項目後,於瀏覽器輸入地址:http://127.0.0.1:8080/index(*根據app.conf的個人配置會有所不同)

package teach_tool

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"goskeleton/app/global/variable"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// FontType 封面处理
type FontType *truetype.Font

type TextType struct {
	Text     string
	Dpi      float64
	FontSize float64
	PosX     int
	PosY     int
	Color    image.Image
	Font     FontType
}

type TextBounds struct {
	Text   string
	Width  float64
	Height float64
}

func FullTemplateImage(templateId int, pos string) string {
	return variable.ConfigCustomYml.GetString("TeachConfig.CoverTemplatePath") + "tpl_" + strconv.Itoa(templateId) + "_" + pos + ".jpg"
}

func GenerateCover(title string, templateId int) (string, error) {
	var fontRegular FontType
	fontRegular, _ = parseFont(variable.BasePath + variable.ConfigCustomYml.GetString("TeachConfig.CoverTitleFont"))

	//创建画布
	im := imaging.New(540, 304, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	//打开模板
	barcodeIm, err := imaging.Open(variable.BasePath + FullTemplateImage(templateId, "cover"))
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("failed to open Cover Template image.  err: %v", err))
		return "", err
	}
	//哈并模板到画布
	draw.Draw(im, barcodeIm.Bounds(), barcodeIm, image.Point{}, draw.Over)
	text := &TextType{Text: title, FontSize: 40, PosX: 35, PosY: 85, Font: fontRegular, Color: image.NewUniform(color.NRGBA{R: 255, G: 255, B: 255, A: 255})}
	textBounds, _ := text.textWidth()
	textLines := textLineFeed(textBounds, 290)
	for i, line := range textLines {
		text.Text = line.Text
		text.PosY = text.PosY + (i * int(line.Height)) + 15
		im, _ = drawText(im, text)
		// 最多只支持写两行文字，多于的省略
		if i >= 1 {
			break
		}
	}

	//imaging.Save(tplIm, path.Join("./results", c.Code+"_front.png"))
	savePath, returnPath := generateSavePath()
	saveFileName := fmt.Sprintf("%d%s", variable.SnowFlake.GetId(), ".jpg")
	_ = imaging.Save(im, savePath+saveFileName)
	return strings.ReplaceAll(returnPath+saveFileName, variable.BasePath, ""), nil
}

func drawText(dstIm image.Image, text *TextType) (*image.NRGBA, error) {
	newIm := image.NewRGBA(dstIm.Bounds())
	draw.Draw(newIm, dstIm.Bounds(), dstIm, image.Point{}, draw.Src)

	c := freetype.NewContext()
	//设置屏幕每英寸的分辨率
	if text.Dpi > 0 {
		c.SetDPI(text.Dpi)
	}
	//设置用于绘制文本的字体
	c.SetFont(text.Font)
	//以磅为单位设置字体大小 (1像素 = 3/4磅?)
	c.SetFontSize(text.FontSize)
	//设置剪裁矩形以进行绘制
	c.SetClip(newIm.Bounds())
	//设置目标图像
	c.SetDst(newIm)
	//设置绘制操作的源图像（颜色），通常为 image.Uniform
	c.SetSrc(text.Color)
	// 写入文字
	_, err := c.DrawString(text.Text, freetype.Pt(text.PosX, text.PosY+int(c.PointToFixed(text.FontSize)>>6)))
	if err != nil {
		variable.ZapLog.Error(err.Error())
		return nil, err
	}
	return (*image.NRGBA)(newIm), nil
}

// textWidth 计算字体文本宽度
func (t *TextType) textWidth() ([]TextBounds, error) {
	opts := truetype.Options{}
	opts.Size = t.FontSize
	face := truetype.NewFace(t.Font, &opts)
	var bounds = make([]TextBounds, 0)
	for _, x := range t.Text {
		aWidth, ok := face.GlyphAdvance(x)
		if ok != true {
			variable.ZapLog.Error("failed to get text width")
			return nil, errors.New("failed to get text width")
		}

		bounds = append(bounds, TextBounds{
			Text:   string(x),
			Width:  float64(aWidth) / 64,
			Height: t.FontSize,
		})
	}
	return bounds, nil
}

// textLineFeed  根据计算出来的文本宽度来换行
func textLineFeed(bounds []TextBounds, lineMaxWidth float64) []TextBounds {
	var lineBounds = make([]TextBounds, 0)
	var lineWidth, lineHeight float64
	var textString string
	for _, v := range bounds {
		lineHeight = v.Height
		if lineWidth+v.Width > lineMaxWidth {

			//换行
			lineBounds = append(lineBounds, TextBounds{
				Text:   textString,
				Width:  lineWidth,
				Height: lineHeight,
			})
			//初始化
			textString = v.Text
			lineWidth = v.Width
		} else {
			textString += v.Text
			lineWidth += v.Width
		}
	}
	if textString != "" {
		lineBounds = append(lineBounds, TextBounds{
			Text:   textString,
			Width:  lineWidth,
			Height: lineHeight,
		})
	}
	return lineBounds
}

/**
 * parseFont 解析字体文件
 */
func parseFont(fontPath string) (FontType, error) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalf("failed to generate QR code: %v", err)
		return nil, err
	}
	return font, nil
}

func generateSavePath() (string, string) {
	savePathPre := variable.BasePath + variable.ConfigCustomYml.GetString("TeachConfig.CoverSavePath")
	returnPath := variable.BasePath + variable.ConfigCustomYml.GetString("TeachConfig.CoverReturnPath")
	curYearMonth := time.Now().Format("2006_01")
	newSavePathPre := savePathPre + curYearMonth
	newReturnPathPre := returnPath + curYearMonth
	// 相关路径不存在，创建目录
	if _, err := os.Stat(newSavePathPre); err != nil {
		if err = os.MkdirAll(newSavePathPre, os.ModePerm); err != nil {
			variable.ZapLog.Error("文件上传创建目录出错" + err.Error())
			return "", ""
		}
	}
	return newSavePathPre + "/", newReturnPathPre + "/"
}

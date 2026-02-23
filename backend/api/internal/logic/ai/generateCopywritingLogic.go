package ai

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCopywritingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCopywritingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCopywritingLogic {
	return &GenerateCopywritingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCopywritingLogic) GenerateCopywriting(req *types.GenerateCopywritingReq) (resp *types.GenerateCopywritingResp, err error) {
	style := req.Style
	if style == "" {
		style = "professional"
	}

	length := req.Length
	if length == "" {
		length = "medium"
	}

	content := l.generateTemplateCopywriting(req.Topic, style, length)

	return &types.GenerateCopywritingResp{
		Content: content,
	}, nil
}

func (l *GenerateCopywritingLogic) generateTemplateCopywriting(topic, style, length string) string {
	var builder strings.Builder

	switch style {
	case "casual":
		builder.WriteString("🎉 ")
		builder.WriteString(topic)
		builder.WriteString(" 来啦！")
	case "urgent":
		builder.WriteString("⚡ 限时特惠：")
		builder.WriteString(topic)
		builder.WriteString("！")
	case "emotional":
		builder.WriteString("❤️ ")
		builder.WriteString(topic)
		builder.WriteString("，让生活更美好！")
	default:
		builder.WriteString("📢 ")
		builder.WriteString(topic)
		builder.WriteString(" 正式上线！")
	}

	builder.WriteString(" 即刻参与，享受专属优惠。")

	switch length {
	case "short":
		break
	case "long":
		builder.WriteString(" 无论是品质还是服务，我们都追求卓越。现在加入，更有限时福利等你来拿！数量有限，先到先得。快来体验吧！")
	default:
		builder.WriteString(" 名额有限，赶快行动！")
	}

	return builder.String()
}

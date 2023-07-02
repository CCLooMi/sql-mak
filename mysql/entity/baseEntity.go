package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type BaseEntity interface {
	TableName() string
}

type IdEntity struct {
	BaseEntity
	ID []byte `gorm:"type:binary(16); primaryKey; not null; comment:'主键ID'"`
}

type TimeEntity struct {
	InsertedAt time.Time `gorm:"not null; comment:'插入时间'"`
	UpdatedAt  time.Time `gorm:"not null; comment:'更新时间'"`
}

type Account struct {
	IdEntity
	TimeEntity
	UserID  []byte          `gorm:"type:binary(16); comment:'用户ID'"`
	Balance decimal.Decimal `gorm:"type:decimal(19,2); comment:'资金'"`
}

func (*Account) TableName() string {
	return "accounts"
}

type Category struct {
	IdEntity
	TimeEntity
	Name        string `gorm:"type:varchar(50); comment:'分类名称'"`
	Description string `gorm:"type:varchar(255); comment:'分类描述'"`
	Order       int    `gorm:"comment:'分类排序'"`
}

func (*Category) TableName() string {
	return "categories"
}

type Comment struct {
	IdEntity
	TimeEntity
	Content  string `gorm:"type:text; comment:'评论内容'"`
	Rating   int    `gorm:"comment:'评分'"`
	UserID   []byte `gorm:"type:binary(16); comment:'用户ID'"`
	TargetID []byte `gorm:"type:binary(16); comment:'目标ID'"`
	RootID   []byte `gorm:"type:binary(16); comment:'根ID'"`
}

func (*Comment) TableName() string {
	return "comments"
}

type Organization struct {
	IdEntity
	TimeEntity
	Name        string `gorm:"type:varchar(255); not null; comment:'组织名称'"`
	Description string `gorm:"type:varchar(255); comment:'组织描述'"`
}

func (*Organization) TableName() string {
	return "organizations"
}

type Permission struct {
	IdEntity
	TimeEntity
	Name        string `gorm:"type:varchar(255); not null; comment:'权限名称'"`
	Descriptor  string `gorm:"type:varchar(255); not null; comment:'权限描述'"`
	Type        string `gorm:"type:varchar(255); not null; comment:'权限类型'"`
	Description string `gorm:"type:varchar(255); comment:'权限描述'"`
}

func (*Permission) TableName() string {
	return "permissions"
}

type PurchasedWpp struct {
	IdEntity
	TimeEntity
	UserID       []byte          `gorm:"type:binary(16); comment:'用户ID'"`
	WppID        []byte          `gorm:"type:binary(16); comment:'应用ID'"`
	Price        decimal.Decimal `gorm:"type:decimal(10, 0); comment:'购买价格'"`
	PurchaseTime time.Time       `gorm:"comment:'购买时间'"`
}

func (*PurchasedWpp) TableName() string {
	return "purchased_wpps"
}

type RolePermission struct {
	IdEntity
	TimeEntity
	RoleID       []byte `gorm:"type:binary(16); comment:'角色ID'"`
	PermissionID []byte `gorm:"type:binary(16); comment:'权限ID'"`
}

func (*RolePermission) TableName() string {
	return "role_permissions"
}

type Role struct {
	IdEntity
	TimeEntity
	Name        string `gorm:"type:varchar(255); not null; comment:'角色名称'"`
	Description string `gorm:"type:varchar(255); comment:'角色描述'"`
}

func (*Role) TableName() string {
	return "roles"
}

type SchemaMigration struct {
	Version    int64      `gorm:"primaryKey; not null; comment:'版本号'"`
	InsertedAt *time.Time `gorm:"comment:'插入时间'"`
}

func (*SchemaMigration) TableName() string {
	return "schema_migrations"
}

type TMessage struct {
	IdEntity
	TimeEntity
	RoomID  string `gorm:"type:varchar(255); comment:'房间ID'"`
	Name    string `gorm:"type:varchar(255); comment:'名称'"`
	Message string `gorm:"type:varchar(255); comment:'消息内容'"`
}

func (*TMessage) TableName() string {
	return "t_messages"
}

type Upload struct {
	IdEntity
	TimeEntity
	FileID   []byte `gorm:"type:varbinary(32); comment:'文件ID'"`
	FileName string `gorm:"type:varchar(255); comment:'文件名称'"`
	FileType string `gorm:"type:varchar(255); comment:'文件类型'"`
	FileSize int64  `gorm:"type:bigint(20); comment:'文件大小'"`
	BizID    []byte `gorm:"type:binary(16); comment:'业务ID'"`
	BizType  string `gorm:"type:varchar(255); comment:'业务类型'"`
}

func (*Upload) TableName() string {
	return "uploads"
}

type UserOrganization struct {
	IdEntity
	TimeEntity
	UserID         []byte `gorm:"type:binary(16); comment:'用户ID'"`
	OrganizationID []byte `gorm:"type:binary(16); comment:'组织ID'"`
}

func (*UserOrganization) TableName() string {
	return "user_organizations"
}

type UserRole struct {
	IdEntity
	TimeEntity
	UserID []byte `gorm:"type:binary(16); comment:'用户ID'"`
	RoleID []byte `gorm:"type:binary(16); comment:'角色ID'"`
}

func (*UserRole) TableName() string {
	return "user_roles"
}

type User struct {
	IdEntity
	TimeEntity
	Username string `gorm:"type:varchar(255); comment:'用户名'"`
	Password []byte `gorm:"type:varbinary(32); comment:'用户密码'"`
}

func (*User) TableName() string {
	return "users"
}

type WppCategory struct {
	IdEntity
	TimeEntity
	WppID      []byte `gorm:"type:binary(16); comment:'应用ID'"`
	CategoryID []byte `gorm:"type:binary(16); comment:'分类ID'"`
}

func (*WppCategory) TableName() string {
	return "wpp_categories"
}

type Wpp struct {
	IdEntity
	TimeEntity
	Name        string `gorm:"type:varchar(64); comment:'应用名称'"`
	Description string `gorm:"type:text; comment:'描述'"`
	Version     string `gorm:"type:varchar(255); comment:'版本号'"`
	DeveloperID []byte `gorm:"type:binary(16); comment:'开发者ID'"`
	FileID      []byte `gorm:"type:varbinary(32); comment:'文件ID'"`
}

func (*Wpp) TableName() string {
	return "wpps"
}

package cat

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Cat struct {
	ID uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`

	CatCode             string `gorm:"type:varchar(64);uniqueIndex;not null;comment:猫咪编号" json:"cat_code"`
	Name                string `gorm:"type:varchar(64);not null;comment:猫咪名称" json:"name"`
	Gender              string `gorm:"type:varchar(20);not null;default:'unknown';comment:性别 male/female/unknown" json:"gender"`
	BodySize            string `gorm:"type:varchar(20);default:'medium';comment:体型 small/medium/large" json:"body_size"`
	AgeStage            string `gorm:"type:varchar(30);default:'young';comment:年龄阶段 kitten/young/adult/old" json:"age_stage"`
	SterilizationStatus string `gorm:"type:varchar(30);default:'unsterilized';comment:绝育状态 sterilized/unsterilized" json:"sterilization_status"`
	AdoptionStatus      string `gorm:"type:varchar(30);default:'pending';comment:领养状态 adopted/pending/unavailable" json:"adoption_status"`
	Avatar              string `gorm:"type:varchar(255);comment:猫咪头像" json:"avatar"`
	Description         string `gorm:"type:text;comment:简介" json:"description"`

	DiscoveryAddress        string  `gorm:"type:varchar(255);comment:发现地址" json:"discovery_address"`
	Longitude               float64 `gorm:"type:decimal(10,6);default:0;comment:经度" json:"longitude"`
	Latitude                float64 `gorm:"type:decimal(10,6);default:0;comment:纬度" json:"latitude"`
	IsVaccinated            bool    `gorm:"default:false;comment:是否已接种疫苗" json:"is_vaccinated"`
	IsHealthy               bool    `gorm:"default:false;comment:是否健康" json:"is_healthy"`
	NeedMedicalIntervention bool    `gorm:"default:true;comment:是否需要医疗干预" json:"need_medical_intervention"`

	CreatedAt time.Time             `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano;comment:删除时间" json:"deleted_at,omitempty"`
}

func (Cat) TableName() string {
	return "cats"
}

// todo 可能需要一个单独的 tag 服务

type Tag struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	Name     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:标签名称" json:"name"`
	Type     string `gorm:"type:varchar(32);index;comment:标签类型" json:"type"`
	Theme    string `gorm:"type:varchar(32);comment:主题色 success/warning/danger/info" json:"theme"`
	Priority int32  `gorm:"default:0;comment:优先级，越大越靠前" json:"priority"`

	CreatedAt time.Time             `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano;comment:删除时间" json:"deleted_at,omitempty"`
}

type TagRelation struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	TagID      uint64    `gorm:"index;not null" json:"tag_id"`
	TargetID   uint64    `gorm:"index;not null" json:"target_id"`
	TargetType string    `gorm:"type:varchar(32);index;not null;comment:cat/post/task/user" json:"target_type"`
	CreatedAt  time.Time `json:"created_at"`
}

type CatImage struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	CatID       uint64 `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`
	ImageURL    string `gorm:"type:varchar(500);not null;comment:图片URL" json:"image_url"`
	ImageType   string `gorm:"type:varchar(32);default:'normal';comment:图片类型 avatar/rescue/medical/adoption/post" json:"image_type"`
	Sort        int32  `gorm:"default:0;comment:排序，越大越靠前" json:"sort"`
	IsCover     bool   `gorm:"default:false;comment:是否封面图" json:"is_cover"`
	Description string `gorm:"type:varchar(255);comment:图片说明" json:"description"`
	UploaderID  uint64 `gorm:"default:0;comment:上传人ID" json:"uploader_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

// CatRescueRecord 救助历程表
type CatRescueRecord struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID        uint64    `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`
	RescueTime   time.Time `gorm:"comment:救助时间" json:"rescue_time"`
	RescuerName  string    `gorm:"type:varchar(64);comment:救助者姓名" json:"rescuer_name"`
	RescueStatus string    `gorm:"type:varchar(50);comment:救助状态" json:"rescue_status"`
	Description  string    `gorm:"type:text;comment:救助简介" json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (CatRescueRecord) TableName() string {
	return "cat_rescue_records"
}

type CatMedicalRecord struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID        uint64    `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`
	MedicalDate  time.Time `gorm:"comment:医疗日期" json:"medical_date"`
	MedicalType  string    `gorm:"type:varchar(50);comment:类型 vaccine/sterilization/check/treatment" json:"medical_type"`
	Content      string    `gorm:"type:text;comment:医疗内容" json:"content"`
	OperatorName string    `gorm:"type:varchar(64);comment:操作员" json:"operator_name"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano;comment:删除时间" json:"deleted_at,omitempty"`
}

func (CatMedicalRecord) TableName() string {
	return "cat_medical_records"
}

// CatTask  猫咪任务表
type CatTask struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID           uint64    `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`
	Title           string    `gorm:"type:varchar(128);not null;comment:任务标题" json:"title"`
	TaskType        string    `gorm:"type:varchar(50);comment:任务类型 feeding/rescue/vaccine/sterilization/adoption" json:"task_type"`
	UrgencyLevel    string    `gorm:"type:varchar(30);default:'normal';comment:紧急程度 high/urgent/normal" json:"urgency_level"`
	DifficultyLevel int32     `gorm:"default:1;comment:难度等级 1-5" json:"difficulty_level"`
	RewardPoints    int32     `gorm:"default:0;comment:积分奖励" json:"reward_points"`
	Status          string    `gorm:"type:varchar(30);default:'pending';comment:状态 pending/processing/completed" json:"status"`
	Summary         string    `gorm:"type:text;comment:任务简述" json:"summary"`
	Detail          string    `gorm:"type:longtext;comment:任务详细说明" json:"detail"`
	Deadline        time.Time `gorm:"comment:截止日期" json:"deadline"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano;comment:删除时间" json:"deleted_at,omitempty"`
}

func (CatTask) TableName() string {
	return "cat_tasks"
}

// CatAdoption 猫咪领养信息表
type CatAdoption struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID               uint64 `gorm:"uniqueIndex;not null;comment:猫咪ID" json:"cat_id"`
	AdoptionRequirement string `gorm:"type:text;comment:领养要求" json:"adoption_requirement"`
	RequiredPoints      int32  `gorm:"default:0;comment:所需积分" json:"required_points"`
	Status              string `gorm:"type:varchar(30);default:'pending';comment:状态 pending/adopted/unavailable" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt
}

func (CatAdoption) TableName() string {
	return "cat_adoptions"
}

/*
 * =========================================================
 * 创建猫咪档案申请表
 * =========================================================
 */

type CatCreateApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`

	// 基础信息
	Name        string `gorm:"type:varchar(64);not null;comment:猫咪名称" json:"name"`
	Gender      string `gorm:"type:varchar(20);default:'unknown';comment:性别" json:"gender"`
	BodySize    string `gorm:"type:varchar(20);default:'medium';comment:体型" json:"body_size"`
	AgeStage    string `gorm:"type:varchar(30);default:'young';comment:年龄阶段" json:"age_stage"`
	Description string `gorm:"type:text;comment:简介" json:"description"`

	// 位置信息
	DiscoveryAddress string  `gorm:"type:varchar(255);comment:发现地址" json:"discovery_address"`
	Longitude        float64 `gorm:"type:decimal(10,6);default:0" json:"longitude"`
	Latitude         float64 `gorm:"type:decimal(10,6);default:0" json:"latitude"`
	ImageURLs        string  `gorm:"type:json;comment:图片URLs JSON" json:"image_urls"`

	// 申请人信息（来自 user-service）
	ApplicantUserID uint64 `gorm:"index;not null;comment:申请人ID" json:"applicant_user_id"`

	// 审核信息
	Status       string `gorm:"type:varchar(30);default:'pending';index;comment:pending/approved/rejected" json:"status"`
	RejectReason string `gorm:"type:text;comment:驳回原因" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0;comment:审核人ID" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatCreateApply) TableName() string {
	return "cat_create_applies"
}

/*
 * =========================================================
 * 修改猫咪档案申请表
 * =========================================================
 */

type CatUpdateApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID uint64 `gorm:"index;not null;comment:猫咪ID" json:"cat_id"`

	Name        string `gorm:"type:varchar(64)" json:"name"`
	Gender      string `gorm:"type:varchar(20)" json:"gender"`
	BodySize    string `gorm:"type:varchar(20)" json:"body_size"`
	AgeStage    string `gorm:"type:varchar(30)" json:"age_stage"`
	Description string `gorm:"type:text" json:"description"`

	DiscoveryAddress string  `gorm:"type:varchar(255)" json:"discovery_address"`
	Longitude        float64 `gorm:"type:decimal(10,6);default:0" json:"longitude"`
	Latitude         float64 `gorm:"type:decimal(10,6);default:0" json:"latitude"`

	ImageURLs string `gorm:"type:json;comment:图片JSON" json:"image_urls"`

	ChangeReason string `gorm:"type:text;comment:修改原因" json:"change_reason"`

	ApplicantUserID uint64 `gorm:"index;not null" json:"applicant_user_id"`

	Status       string `gorm:"type:varchar(30);default:'pending';index" json:"status"`
	RejectReason string `gorm:"type:text" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatUpdateApply) TableName() string {
	return "cat_update_applies"
}

/*
 * =========================================================
 * 医疗记录申请表
 * =========================================================
 */

type CatMedicalApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID uint64 `gorm:"index;not null" json:"cat_id"`

	MedicalType string `gorm:"type:varchar(50);comment:vaccine/sterilization/check/treatment" json:"medical_type"`
	Content     string `gorm:"type:text" json:"content"`

	ApplicantUserID uint64 `gorm:"index;not null" json:"applicant_user_id"`

	Status       string `gorm:"type:varchar(30);default:'pending';index" json:"status"`
	RejectReason string `gorm:"type:text" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatMedicalApply) TableName() string {
	return "cat_medical_applies"
}

/*
 * =========================================================
 * 救助记录申请表
 * =========================================================
 */

type CatRescueApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID uint64 `gorm:"index;not null" json:"cat_id"`

	RescueTime   time.Time `json:"rescue_time"`
	RescuerName  string    `gorm:"type:varchar(64)" json:"rescuer_name"`
	RescueStatus string    `gorm:"type:varchar(50)" json:"rescue_status"`
	Description  string    `gorm:"type:text" json:"description"`

	ApplicantUserID uint64 `gorm:"index;not null" json:"applicant_user_id"`

	Status       string `gorm:"type:varchar(30);default:'pending';index" json:"status"`
	RejectReason string `gorm:"type:text" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatRescueApply) TableName() string {
	return "cat_rescue_applies"
}

/*
 * =========================================================
 * 任务创建申请表
 * =========================================================
 */

type CatTaskApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CatID uint64 `gorm:"index;not null" json:"cat_id"`

	Title    string `gorm:"type:varchar(128);not null" json:"title"`
	TaskType string `gorm:"type:varchar(50)" json:"task_type"`

	Summary string `gorm:"type:text" json:"summary"`
	Detail  string `gorm:"type:longtext" json:"detail"`

	Deadline time.Time `json:"deadline"`

	ApplicantUserID uint64 `gorm:"index;not null" json:"applicant_user_id"`

	// 审核时由管理员填写
	UrgencyLevel    string `gorm:"type:varchar(30);default:'normal'" json:"urgency_level"`
	DifficultyLevel int32  `gorm:"default:1" json:"difficulty_level"`
	RewardPoints    int32  `gorm:"default:0" json:"reward_points"`

	Status       string `gorm:"type:varchar(30);default:'pending';index" json:"status"`
	RejectReason string `gorm:"type:text" json:"reject_reason"`
	ReviewerID   uint64 `gorm:"default:0" json:"reviewer_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatTaskApply) TableName() string {
	return "cat_task_applies"
}

/*
 * =========================================================
 * 任务认领关系表（非常重要）
 * =========================================================
 */

type CatTaskClaim struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	TaskID uint64 `gorm:"index;not null;comment:任务ID" json:"task_id"`
	UserID uint64 `gorm:"index;not null;comment:认领用户ID" json:"user_id"`

	Status string `gorm:"type:varchar(30);default:'claimed';comment:claimed/completed/abandoned" json:"status"`

	ClaimTime  time.Time `json:"claim_time"`
	FinishTime time.Time `json:"finish_time"`

	AbandonReason string `gorm:"type:text" json:"abandon_reason"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatTaskClaim) TableName() string {
	return "cat_task_claims"
}

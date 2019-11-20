package gqlgen

import (
	"context"
	"demo13/prisma"
	"demo13/servants/AttendsDefaultDataManager"
	// "log"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{
	Prisma *prisma.Client
}

func (r *Resolver) Attend() AttendResolver {
	return &attendResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Staff() StaffResolver {
	return &staffResolver{r}
}

type attendResolver struct{ *Resolver }

func (r *attendResolver) StaffInfo(ctx context.Context, obj *prisma.Attend) (*prisma.Staff, error) {
	staff, err := r.Prisma.Attend(prisma.AttendWhereUniqueInput{
		ID: &obj.ID,
	}).StaffInfo().Exec(ctx)
	return staff, err
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateStaff(ctx context.Context, name string, age int, profileImagePath *string) (*prisma.Staff, error) {
	addm := AttendsDefaultDataManager.Summon()
	staff, err := r.Prisma.CreateStaff(prisma.StaffCreateInput{
		Name:             &name,
		Age:             prisma.Int32(int32(age)),
		ProfileImagePath: profileImagePath,
		Attends: &prisma.AttendCreateManyWithoutStaffInfoInput{
			Create: addm.ServeData(),
		},
	}).Exec(ctx)

	return staff, err
}
func (r *mutationResolver) UpdateStaffProfile(ctx context.Context, id string, name *string, age *int, profileImagePath *string) (*prisma.Staff, error) {
	staff, _ := r.Prisma.Staff(prisma.StaffWhereUniqueInput{ ID: &id }).Exec(ctx)

	if name == nil { name = staff.Name }
	if age == nil { a := *staff.Age; a2 := int(a); age = &a2 }
	if profileImagePath == nil { profileImagePath = staff.ProfileImagePath }
	
	re, err := r.Prisma.UpdateStaff(prisma.StaffUpdateParams{
		Where: prisma.StaffWhereUniqueInput{ ID: &id },
		Data: prisma.StaffUpdateInput{
			Name: name,
			Age: prisma.Int32(int32(*age)),
			ProfileImagePath: profileImagePath,
		},
	}).Exec(ctx)
	return re, err
}
func (r *mutationResolver) UpdateStaffAttend(ctx context.Context, staffID string, attendID string, input UpdateStaffAttendInput) (*prisma.Attend, error) {
	attend, _ := r.Prisma.Attend(prisma.AttendWhereUniqueInput{ ID: &attendID }).Exec(ctx)

	if input.IsAttend == nil { input.IsAttend = &attend.IsAttend }
	if input.InTimeIndex == nil { i := attend.InTimeIndex; i2 := int(i); input.InTimeIndex = &i2 }
	if input.OutTimeIndex == nil { o := attend.OutTimeIndex; o2 := int(o); input.OutTimeIndex = &o2 }
	
	_, err := r.Prisma.UpdateManyAttends(prisma.AttendUpdateManyParams{
		Where: &prisma.AttendWhereInput{
			And: []prisma.AttendWhereInput{
				prisma.AttendWhereInput{ ID: &attendID },
				prisma.AttendWhereInput{
					StaffInfo: &prisma.StaffWhereInput{ ID: &staffID },
				},
			},
		},
		Data: prisma.AttendUpdateManyMutationInput{
			IsAttend: input.IsAttend,
			InTimeIndex: prisma.Int32(int32(*input.InTimeIndex)),
			OutTimeIndex: prisma.Int32(int32(*input.OutTimeIndex)),
		},
	}).Exec(ctx)

	if err != nil { panic(err) }

	re, _ := r.Prisma.Attend(prisma.AttendWhereUniqueInput{ ID: &attendID }).Exec(ctx)

	return re, err
}
func (r *mutationResolver) DeleteStaff(ctx context.Context, id string) (*prisma.Staff, error) {
	staff, err := r.Prisma.DeleteStaff(prisma.StaffWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	return staff, err
}
func (r *mutationResolver) CronUpdateAttend(ctx context.Context, apiKey string) (*prisma.Attend, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Staffs(ctx context.Context) ([]*prisma.Staff, error) {
	re := []*prisma.Staff{}
	staffs, err := r.Prisma.Staffs(nil).Exec(ctx)

	for _, staff := range staffs {
		re = append(re, &prisma.Staff{
			ID:               staff.ID,
			Name:             staff.Name,
			Age:              staff.Age,
			ProfileImagePath: staff.ProfileImagePath,
			CreatedAt:        staff.CreatedAt,
			UpdatedAt:        staff.UpdatedAt,
		})
	}

	return re, err
}
func (r *queryResolver) Staff(ctx context.Context, id string) (*prisma.Staff, error) {
	staff, err := r.Prisma.Staff(prisma.StaffWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	return staff, err
}
func (r *queryResolver) Attend(ctx context.Context, staffID string, attendID string) (*prisma.Attend, error) {
	panic("not implemented")
}

type staffResolver struct{ *Resolver }

func (r *staffResolver) Attends(ctx context.Context, obj *prisma.Staff) ([]*prisma.Attend, error) {
	re := []*prisma.Attend{}

	attends, err := r.Prisma.Staff(prisma.StaffWhereUniqueInput{
		ID: &obj.ID,
	}).Attends(nil).Exec(ctx)

	for _, attend := range attends {
		re = append(re, &prisma.Attend{
			ID:            attend.ID,
			DateStartTime: attend.DateStartTime,
			IsAttend:      attend.IsAttend,
			InTimeIndex:   attend.InTimeIndex,
			OutTimeIndex:  attend.OutTimeIndex,
			CreatedAt:     attend.CreatedAt,
			UpdatedAt:     attend.UpdatedAt,
		})
	}

	return re, err
}

package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetAttackTechniques retrieves all MITRE ATT&CK techniques.
func (c *Client) GetAttackTechniques(ctx context.Context) ([]*types.Attack, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Attack []struct {
			ID     int    `graphql:"id"`
			TNum   string `graphql:"t_num"`
			Name   string `graphql:"name"`
			OS     string `graphql:"os"`
			Tactic string `graphql:"tactic"`
		} `graphql:"attack(order_by: {t_num: asc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetAttackTechniques", err, "failed to query MITRE ATT&CK techniques")
	}

	attacks := make([]*types.Attack, len(query.Attack))
	for i, a := range query.Attack {
		attacks[i] = &types.Attack{
			ID:     a.ID,
			TNum:   a.TNum,
			Name:   a.Name,
			OS:     a.OS,
			Tactic: a.Tactic,
		}
	}

	return attacks, nil
}

// GetAttackTechniqueByID retrieves a specific MITRE ATT&CK technique by ID.
func (c *Client) GetAttackTechniqueByID(ctx context.Context, attackID int) (*types.Attack, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if attackID == 0 {
		return nil, WrapError("GetAttackTechniqueByID", ErrInvalidInput, "attack ID is required")
	}

	var query struct {
		Attack []struct {
			ID     int    `graphql:"id"`
			TNum   string `graphql:"t_num"`
			Name   string `graphql:"name"`
			OS     string `graphql:"os"`
			Tactic string `graphql:"tactic"`
		} `graphql:"attack(where: {id: {_eq: $attack_id}})"`
	}

	variables := map[string]interface{}{
		"attack_id": attackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetAttackTechniqueByID", err, "failed to query MITRE ATT&CK technique")
	}

	if len(query.Attack) == 0 {
		return nil, WrapError("GetAttackTechniqueByID", ErrNotFound, "ATT&CK technique not found")
	}

	a := query.Attack[0]
	return &types.Attack{
		ID:     a.ID,
		TNum:   a.TNum,
		Name:   a.Name,
		OS:     a.OS,
		Tactic: a.Tactic,
	}, nil
}

// GetAttackTechniqueByTNum retrieves a MITRE ATT&CK technique by technique number (e.g., "T1003").
func (c *Client) GetAttackTechniqueByTNum(ctx context.Context, tNum string) (*types.Attack, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if tNum == "" {
		return nil, WrapError("GetAttackTechniqueByTNum", ErrInvalidInput, "technique number is required")
	}

	var query struct {
		Attack []struct {
			ID     int    `graphql:"id"`
			TNum   string `graphql:"t_num"`
			Name   string `graphql:"name"`
			OS     string `graphql:"os"`
			Tactic string `graphql:"tactic"`
		} `graphql:"attack(where: {t_num: {_eq: $t_num}})"`
	}

	variables := map[string]interface{}{
		"t_num": tNum,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetAttackTechniqueByTNum", err, "failed to query MITRE ATT&CK technique")
	}

	if len(query.Attack) == 0 {
		return nil, WrapError("GetAttackTechniqueByTNum", ErrNotFound, "ATT&CK technique not found")
	}

	a := query.Attack[0]
	return &types.Attack{
		ID:     a.ID,
		TNum:   a.TNum,
		Name:   a.Name,
		OS:     a.OS,
		Tactic: a.Tactic,
	}, nil
}

// GetAttackByTask retrieves all MITRE ATT&CK tags for a specific task.
func (c *Client) GetAttackByTask(ctx context.Context, taskID int) ([]*types.AttackTask, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if taskID == 0 {
		return nil, WrapError("GetAttackByTask", ErrInvalidInput, "task ID is required")
	}

	var query struct {
		AttackTask []struct {
			ID       int `graphql:"id"`
			AttackID int `graphql:"attack_id"`
			TaskID   int `graphql:"task_id"`
		} `graphql:"attacktask(where: {task_id: {_eq: $task_id}})"`
	}

	variables := map[string]interface{}{
		"task_id": taskID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetAttackByTask", err, "failed to query MITRE ATT&CK tags for task")
	}

	attackTasks := make([]*types.AttackTask, len(query.AttackTask))
	for i, at := range query.AttackTask {
		attackTasks[i] = &types.AttackTask{
			ID:       at.ID,
			AttackID: at.AttackID,
			TaskID:   at.TaskID,
		}
	}

	return attackTasks, nil
}

// GetAttackByCommand retrieves all MITRE ATT&CK tags for a specific command.
func (c *Client) GetAttackByCommand(ctx context.Context, commandID int) ([]*types.AttackCommand, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if commandID == 0 {
		return nil, WrapError("GetAttackByCommand", ErrInvalidInput, "command ID is required")
	}

	var query struct {
		AttackCommand []struct {
			ID        int `graphql:"id"`
			AttackID  int `graphql:"attack_id"`
			CommandID int `graphql:"command_id"`
		} `graphql:"attackcommand(where: {command_id: {_eq: $command_id}})"`
	}

	variables := map[string]interface{}{
		"command_id": commandID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetAttackByCommand", err, "failed to query MITRE ATT&CK tags for command")
	}

	attackCommands := make([]*types.AttackCommand, len(query.AttackCommand))
	for i, ac := range query.AttackCommand {
		attackCommands[i] = &types.AttackCommand{
			ID:        ac.ID,
			AttackID:  ac.AttackID,
			CommandID: ac.CommandID,
		}
	}

	return attackCommands, nil
}

// GetAttacksByOperation retrieves all MITRE ATT&CK techniques used in an operation.
func (c *Client) GetAttacksByOperation(ctx context.Context, operationID int) ([]*types.Attack, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetAttacksByOperation", ErrInvalidInput, "operation ID is required")
	}

	// Query attacktask table joined with attack and task tables
	var query struct {
		AttackTask []struct {
			Attack struct {
				ID     int    `graphql:"id"`
				TNum   string `graphql:"t_num"`
				Name   string `graphql:"name"`
				OS     string `graphql:"os"`
				Tactic string `graphql:"tactic"`
			} `graphql:"attack"`
		} `graphql:"attacktask(where: {task: {operation_id: {_eq: $operation_id}}}, order_by: {attack: {t_num: asc}}, distinct_on: attack_id)"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetAttacksByOperation", err, "failed to query MITRE ATT&CK techniques for operation")
	}

	attacks := make([]*types.Attack, len(query.AttackTask))
	for i, at := range query.AttackTask {
		attacks[i] = &types.Attack{
			ID:     at.Attack.ID,
			TNum:   at.Attack.TNum,
			Name:   at.Attack.Name,
			OS:     at.Attack.OS,
			Tactic: at.Attack.Tactic,
		}
	}

	return attacks, nil
}

// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"os"
	"testing"

	"github.com/alexellis/derek/types"
)

var commandTriggers = []string{commandTriggerDefault, commandTriggerSlash}

func Test_getCommandTrigger(t *testing.T) {
	const (
		envVar       = "use_slash_trigger"
		errorMessage = "expected trigger to be: %s, got: %s"
	)
	var trigger string

	// test default trigger
	os.Unsetenv(envVar)
	trigger = getCommandTrigger()
	if trigger != commandTriggerDefault {
		t.Errorf(errorMessage, commandTriggerDefault, trigger)
	}

	// test slash trigger
	os.Setenv(envVar, "true")
	trigger = getCommandTrigger()
	if trigger != commandTriggerSlash {
		t.Errorf(errorMessage, commandTriggerSlash, trigger)
	}
}

func Test_Parsing_OpenClose(t *testing.T) {

	var actionOptions = []struct {
		title          string
		body           string
		expectedAction string
	}{
		{
			title:          "Correct reopen command",
			body:           "reopen",
			expectedAction: "reopen",
		},
		{ //this case replaces Test_Parsing_Close
			title:          "Correct close command",
			body:           "close",
			expectedAction: "close",
		},
		{
			title:          "invalid command",
			body:           "dance",
			expectedAction: "",
		},
		{
			title:          "Longer reopen command",
			body:           "reopen: ",
			expectedAction: "reopen",
		},
		{
			title:          "Longer close command",
			body:           "close: ",
			expectedAction: "close",
		},
	}

	for _, test := range actionOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range commandTriggers {
				action := parse(trigger+test.body, trigger)
				if action.Type != test.expectedAction {
					t.Errorf("Action - want: %s, got %s", test.expectedAction, action.Type)
				}
			}
		})
	}
}

func Test_Parsing_Labels(t *testing.T) {

	var labelOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{ //this case replaces Test_Parsing_AddLabel
			title:        "Add label of demo",
			body:         "add label: demo",
			expectedType: "AddLabel",
			expectedVal:  "demo",
		},
		{
			title:        "Remove label of demo",
			body:         "remove label: demo",
			expectedType: "RemoveLabel",
			expectedVal:  "demo",
		},
		{
			title:        "Invalid label action",
			body:         "peel label: demo",
			expectedType: "",
			expectedVal:  "",
		},
	}

	for _, test := range labelOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range commandTriggers {
				action := parse(trigger+test.body, trigger)
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("Action - wanted: %s, got %s\nLabel - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_Parsing_Assignments(t *testing.T) {

	var assignmentOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{
			title:        "Assign to burt",
			body:         "assign: burt",
			expectedType: assignConstant,
			expectedVal:  "burt",
		},
		{
			title:        "Unassign burt",
			body:         "unassign: burt",
			expectedType: unassignConstant,
			expectedVal:  "burt",
		},
		{
			title:        "Assign to me",
			body:         "assign: me",
			expectedType: assignConstant,
			expectedVal:  "me",
		},
		{
			title:        "Unassign me",
			body:         "unassign: me",
			expectedType: unassignConstant,
			expectedVal:  "me",
		},
		{
			title:        "Invalid assignment action",
			body:         "consign: burt",
			expectedType: "",
			expectedVal:  "",
		},
		{
			title:        "Unassign blank",
			body:         "unassign: ",
			expectedType: "",
			expectedVal:  "",
		},
	}

	for _, test := range assignmentOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range commandTriggers {
				action := parse(trigger+test.body, trigger)
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("Action - wanted: %s, got %s\nMaintainer - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_Parsing_Titles(t *testing.T) {

	var titleOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{
			title:        "Set Title",
			body:         "set title: This is a really great Title!",
			expectedType: setTitleConstant,
			expectedVal:  "This is a really great Title!",
		},
		{
			title:        "Mis-spelling of title",
			body:         "set titel: This is a really great Title!",
			expectedType: "",
			expectedVal:  "",
		},
		{
			title:        "Empty Title",
			body:         "set title: ",
			expectedType: "", //blank because it should fail isValidCommand
			expectedVal:  "",
		},
		{
			title:        "Empty Title (Double Space)",
			body:         "set title:  ",
			expectedType: setTitleConstant,
			expectedVal:  "",
		},
	}

	for _, test := range titleOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range commandTriggers {
				action := parse(trigger+test.body, trigger)
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("\nAction - wanted: %s, got %s\nValue - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_assessState(t *testing.T) {

	var stateOptions = []struct {
		title            string
		requestedAction  string
		currentState     string
		expectedNewState string
		expectedBool     bool
	}{
		{
			title:            "Currently Closed and trying to close",
			requestedAction:  closeConstant,
			currentState:     closedConstant,
			expectedNewState: "",
			expectedBool:     false,
		},
		{
			title:            "Currently Open and trying to reopen",
			requestedAction:  reopenConstant,
			currentState:     openConstant,
			expectedNewState: "",
			expectedBool:     false,
		},
		{
			title:            "Currently Closed and trying to open",
			requestedAction:  reopenConstant,
			currentState:     closedConstant,
			expectedNewState: openConstant,
			expectedBool:     true,
		},
		{
			title:            "Currently Open and trying to close",
			requestedAction:  closeConstant,
			currentState:     openConstant,
			expectedNewState: closedConstant,
			expectedBool:     true,
		},
	}

	for _, test := range stateOptions {
		t.Run(test.title, func(t *testing.T) {

			newState, validTransition := checkTransition(test.requestedAction, test.currentState)

			if newState != test.expectedNewState || validTransition != test.expectedBool {
				t.Errorf("\nStates - wanted: %s, got %s\nValidity - wanted: %t, got %t\n", test.expectedNewState, newState, test.expectedBool, validTransition)
			}
		})
	}
}

func Test_validAction(t *testing.T) {

	var stateOptions = []struct {
		title           string
		running         bool
		requestedAction string
		start           string
		stop            string
		expectedBool    bool
	}{
		{
			title:           "Currently unlocked and trying to lock",
			running:         false,
			requestedAction: lockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    true,
		},
		{
			title:           "Currently unlocked and trying to unlock",
			running:         false,
			requestedAction: unlockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    false,
		},
		{
			title:           "Currently locked and trying to lock",
			running:         true,
			requestedAction: lockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    false,
		},
		{
			title:           "Currently locked and trying to unlock",
			running:         true,
			requestedAction: unlockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    true,
		},
	}

	for _, test := range stateOptions {
		t.Run(test.title, func(t *testing.T) {

			isValid := validAction(test.running, test.requestedAction, test.start, test.stop)

			if isValid != test.expectedBool {
				t.Errorf("\nActions - wanted: %t, got %t\n", test.expectedBool, isValid)
			}
		})
	}
}

func Test_findLabel(t *testing.T) {

	var stateOptions = []struct {
		title         string
		currentLabels []types.IssueLabel
		cmdLabel      string
		expectedFound bool
	}{
		{
			title: "Label exists lowercase",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "rod",
			expectedFound: true,
		},
		{
			title: "Label exists case insensitive",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "Rod",
			expectedFound: true,
		},
		{
			title: "Label doesnt exist lowercase",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "derek",
			expectedFound: false,
		},
		{
			title: "Label doesnt exist case insensitive",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "Derek",
			expectedFound: false,
		},
		{
			title:         "no existing labels lowercase",
			currentLabels: nil,
			cmdLabel:      "derek",
			expectedFound: false,
		},
		{title: "Label doesnt exist case insensitive",
			currentLabels: nil,
			cmdLabel:      "Derek",
			expectedFound: false,
		},
	}

	for _, test := range stateOptions {
		t.Run(test.title, func(t *testing.T) {

			labelFound := findLabel(test.currentLabels, test.cmdLabel)

			if labelFound != test.expectedFound {
				t.Errorf("Find Labels(%s) - wanted: %t, got %t\n", test.title, test.expectedFound, labelFound)
			}
		})
	}
}

func Test_Parsing_Milestones(t *testing.T) {

	var milestonesOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{
			title:        "Right set milestone",
			body:         "set milestone: demo",
			expectedType: "SetMilestone",
			expectedVal:  "demo",
		},
		{
			title:        "Right remove milestone",
			body:         "remove milestone: demo",
			expectedType: "RemoveMilestone",
			expectedVal:  "demo",
		},
		{
			title:        "Wrong set milestone",
			body:         "you ok label: demo",
			expectedType: "",
			expectedVal:  "",
		},
	}

	for _, test := range milestonesOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range commandTriggers {
				action := parse(trigger+test.body, trigger)
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("Action - wanted: %s, got %s\nLabel - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_isDcoLabel(t *testing.T) {
	dcoLabel := []struct {
		title        string
		label        string
		expectedBool bool
	}{
		{
			title:        "Counts as no-dco - case insensitivity",
			label:        "NO-DCO",
			expectedBool: true,
		},
		{
			title:        "Normal no-dco case",
			label:        "no-dco",
			expectedBool: true,
		},
		{
			title:        "Counts as no-dco - case insensitivity",
			label:        "No-Dco",
			expectedBool: true,
		},
		{
			title:        "Does not follow no-dco so it counts as normal label",
			label:        "nodco",
			expectedBool: false,
		},
		{
			title:        "Normal label",
			label:        "randomlabel",
			expectedBool: false,
		},
	}

	for _, test := range dcoLabel {
		t.Run(test.label, func(t *testing.T) {
			itsDco := isDcoLabel(test.label)
			if itsDco != test.expectedBool {
				t.Errorf("Wanted `%s` to return: %t but it returned:  %t.", test.label, test.expectedBool, itsDco)
			}
		})
	}
}

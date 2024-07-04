Feature: Animation

  Scenario: Creating a new animation
    Given I have a set of frames
    When I create a new animation with frame duration 0.2, looping enabled, and max loops 3
    Then the animation should have a frame duration of 0.2
    And the animation should loop
    And the animation should have a maximum of 3 loops

  Scenario: Updating the animation
    Given I have a new animation with 2 frames and frame duration 0.1
    When I update the animation with 0.15 seconds
    Then the current frame should be 1

  Scenario: Pausing and resuming the animation
    Given I have a new animation with 2 frames and frame duration 0.1
    And I pause the animation
    When I update the animation with 0.15 seconds
    Then the current frame should be 0
    And I resume the animation
    When I update the animation with 0.15 seconds
    Then the current frame should be 1

  Scenario: Resetting the animation
    Given I have a new animation with 2 frames and frame duration 0.1
    When I update the animation with 0.15 seconds
    And I reset the animation
    Then the current frame should be 0
    And the elapsed time should be 0


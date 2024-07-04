Feature: AudioManager

  Scenario: Loading a sound
    Given I have an AudioManager with a sample rate of 44100
    When I load a sound from "sound.wav" with the name "testSound"
    Then the sound "testSound" should be loaded

  Scenario: Playing a loaded sound
    Given I have an AudioManager with a sample rate of 44100
    And I have loaded a sound from "sound.wav" with the name "testSound"
    When I play the sound "testSound"
    Then the sound should play without errors

  Scenario: Loading and playing background music
    Given I have an AudioManager with a sample rate of 44100
    When I load background music from "bgm.wav"
    And I play the background music
    Then the background music should play without errors

  Scenario: Stopping background music
    Given I have an AudioManager with a sample rate of 44100
    And I have loaded background music from "bgm.wav"
    And I am playing the background music
    When I stop the background music
    Then the background music should stop


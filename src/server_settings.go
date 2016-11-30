package main

// CsgoServerSettings CS:GO Server Settings
type CsgoServerSettings struct {
	AmmoGrenadeLimit                     int     `csgo:"ammo_grenade_limit_default"`
	AmmoGrenadeLimitFlashbang            int     `csgo:"ammo_grenade_limit_flashbang"`
	AmmoGrenadeLimitTotal                int     `csgo:"ammo_grenade_limit_total"`
	BotQuota                             int     `csgo:"bot_quota"`
	CashPlayerBombDefused                int     `csgo:"cash_player_bomb_defused"`
	CashPlayerBombPlanted                int     `csgo:"cash_player_bomb_planted"`
	CashPlayerDamageHostage              int     `csgo:"cash_player_damage_hostage"`
	CashPlayerInteractWithHostage        int     `csgo:"cash_player_interact_with_hostage"`
	CashPlayerKilledEnemyDefault         int     `csgo:"cash_player_killed_enemy_default"`
	CashPlayerKilledEnemyFactor          int     `csgo:"cash_player_killed_enemy_factor"`
	CashPlayerKilledHostage              int     `csgo:"cash_player_killed_hostage"`
	CashPlayerKilledTeammate             int     `csgo:"cash_player_killed_teammate"`
	CashPlayerRescuedHostage             int     `csgo:"cash_player_rescued_hostage"`
	CashTeamEliminationBombMap           int     `csgo:"cash_team_elimination_bomb_map"`
	CashTeamHostageAlive                 int     `csgo:"cash_team_hostage_alive"`
	CashTeamHostageInteraction           int     `csgo:"cash_team_hostage_interaction"`
	CashTeamLoserBonus                   int     `csgo:"cash_team_loser_bonus"`
	CashTeamLoserBonusConsecutiveRounds  int     `csgo:"cash_team_loser_bonus_consecutive_rounds"`
	CashTeamPlantedBombButDefused        int     `csgo:"cash_team_planted_bomb_but_defused"`
	CashTeamRescuedHostage               int     `csgo:"cash_team_rescued_hostage"`
	CashTeamTerroristWinBomb             int     `csgo:"cash_team_terrorist_win_bomb"`
	CashTeamWinByDefusingBomb            int     `csgo:"cash_team_win_by_defusing_bomb"`
	CashTeamWinByHostageRescue           int     `csgo:"cash_team_win_by_hostage_rescue"`
	CashPlayerGetKilled                  int     `csgo:"cash_player_get_killed"`
	CashPlayerRespawnAmount              int     `csgo:"cash_player_respawn_amount"`
	CashTeamEliminationHostageMapCt      int     `csgo:"cash_team_elimination_hostage_map_ct"`
	CashTeamEliminationHostageMapT       int     `csgo:"cash_team_elimination_hostage_map_t"`
	CashTeamWinByTimeRunningOutBomb      int     `csgo:"cash_team_win_by_time_running_out_bomb"`
	CashTeamWinByTimeRunningOutHostage   int     `csgo:"cash_team_win_by_time_running_out_hostage"`
	FfDamageReductionGrenade             float64 `csgo:"ff_damage_reduction_grenade"`      // How much to reduce damage done to teammates by a thrown grenade.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
	FfDamageReductionBullets             float64 `csgo:"ff_damage_reduction_bullets"`      // How much to reduce damage done to teammates when shot.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
	FfDamageReductionOther               float64 `csgo:"ff_damage_reduction_other"`        // How much to reduce damage done to teammates by things other than bullets and grenades.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
	FfDamageReductionGrenadeSelf         float64 `csgo:"ff_damage_reduction_grenade_self"` // How much to damage a player does to himself with his own grenade.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
	MpAfterroundmoney                    int     `csgo:"mp_afterroundmoney"`               // amount of money awared to every player after each round
	MpAutokick                           int     `csgo:"mp_autokick"`                      // Kick idle/team-killing players
	MpAutoteambalance                    int     `csgo:"mp_autoteambalance"`
	MpBuytime                            int     `csgo:"mp_buytime"`                  // How many seconds after round start players can buy items for.
	MpC4timer                            int     `csgo:"mp_c4timer"`                  // How long from when the C4 is armed until it blow"
	MpDeathDropDefuser                   int     `csgo:"mp_death_drop_defuser"`       // Drop defuser on player death
	MpDeathDropGrenade                   int     `csgo:"mp_death_drop_grenade"`       // Which grenade to drop on player death: 0=none, 1=best, 2=current or best
	MpDeathDropGun                       int     `csgo:"mp_death_drop_gun"`           // Which gun to drop on player death: 0=none, 1=best, 2=current or best
	MpDefuserAllocation                  int     `csgo:"mp_defuser_allocation"`       // How to allocate defusers to CTs at start or round: 0=none, 1=random, 2=everyone
	MpDoWarmupPeriod                     int     `csgo:"mp_do_warmup_period"`         // Whether or not to do a warmup period at the start of a match.
	MpForcecamera                        int     `csgo:"mp_forcecamera"`              // Restricts spectator modes for dead players
	MpForcePickTime                      int     `csgo:"mp_force_pick_time"`          // The amount of time a player has on the team screen to make a selection before being auto-teamed
	MpFreeArmor                          int     `csgo:"mp_free_armor"`               // Determines whether armor and helmet are given automatically.
	MpFreezetime                         int     `csgo:"mp_freezetime"`               // How many seconds to keep players frozen when the round starts
	MpFriendlyfire                       int     `csgo:"mp_friendlyfire"`             // Allows team members to injure other members of their team
	MpHalftime                           int     `csgo:"mp_halftime"`                 // Determines whether or not the match has a team-swapping halftime event.
	MpHalftimeDuration                   int     `csgo:"mp_halftime_duration"`        // Number of seconds that halftime lasts
	MpJoinGraceTime                      int     `csgo:"mp_join_grace_time"`          // Number of seconds after round start to allow a player to join a game
	MpLimitteams                         int     `csgo:"mp_limitteams"`               // Max # of players 1 team can have over another (0 disables check)
	MpLogdetail                          int     `csgo:"mp_logdetail"`                // Logs attacks.  Values are: 0=off, 1=enemy, 2=teammate, 3=both)
	MpMatchCanClinch                     int     `csgo:"mp_match_can_clinch"`         // Can a team clinch and end the match by being so far ahead that the other team has no way to catching up
	MpMatchEndRestart                    int     `csgo:"mp_match_end_restart"`        // At the end of the match, perform a restart instead of loading a new map
	MpMaxmoney                           int     `csgo:"mp_maxmoney"`                 // maximum amount of money allowed in a player's account
	MpMaxrounds                          int     `csgo:"mp_maxrounds"`                // max number of rounds to play before server changes maps
	MpMolotovusedelay                    int     `csgo:"mp_molotovusedelay"`          // Number of seconds to delay before the molotov can be used after acquiring it
	MpPlayercashawards                   int     `csgo:"mp_playercashawards"`         // Players can earn money by performing in-game actions
	MpPlayerid                           int     `csgo:"mp_playerid"`                 // Controls what information player see in the status bar: 0 all names; 1 team names; 2 no names
	MpPlayeridDelay                      float64 `csgo:"mp_playerid_delay"`           // Number of seconds to delay showing information in the status bar
	MpPlayeridHold                       float64 `csgo:"mp_playerid_hold"`            // Number of seconds to keep showing old information in the status bar
	MpRoundRestartDelay                  float64 `csgo:"mp_round_restart_delay"`      // Number of seconds to delay before restarting a round after a win
	MpRoundtime                          float64 `csgo:"mp_roundtime"`                // How many minutes each round takes.
	MpRoundtimeDefuse                    float64 `csgo:"mp_roundtime_defuse"`         // How many minutes each round takes on defusal maps.
	MpSolidTeammates                     int     `csgo:"mp_solid_teammates"`          // Determines whether teammates are solid or not.
	MpStartmoney                         int     `csgo:"mp_startmoney"`               // amount of money each player gets when they reset
	MpTeamcashawards                     int     `csgo:"mp_teamcashawards"`           // Teams can earn money by performing in-game actions
	MpTimelimit                          int     `csgo:"mp_timelimit"`                // game time per map in minutes
	MpTkpunish                           int     `csgo:"mp_tkpunish"`                 // Will a TK'er be punished in the next round?  {0=no,  1=yes}
	MpWarmuptime                         int     `csgo:"mp_warmuptime"`               // If true, there will be a warmup period/round at the start of each match to allow
	MpWeaponsAllowMapPlaced              int     `csgo:"mp_weapons_allow_map_placed"` // If this convar is set, when a match starts, the game will not delete weapons placed in the map.
	MpWeaponsAllowZeus                   int     `csgo:"mp_weapons_allow_zeus"`       // Determines whether the Zeus is purchasable or not.
	MpWinPanelDisplayTime                int     `csgo:"mp_win_panel_display_time"`   // The amount of time to show the win panel between matches / halfs
	MpOvertimeEnable                     int     `csgo:"mp_overtime_enable"`
	MpOvertimeMaxrounds                  int     `csgo:"mp_overtime_maxrounds"`
	MpOvertimeStartmoney                 int     `csgo:"mp_overtime_startmoney"`
	RconPassword                         string  `csgo:"rcon_password"`
	TvEnable                             int     `csgo:"tv_enable"`
	TvDelay                              int     `csgo:"tv_delay"`
	SpecFreezeTime                       int     `csgo:"spec_freeze_time"`                // Time spend frozen in observer freeze cam.
	SpecFreezePanelExtendedTime          int     `csgo:"spec_freeze_panel_extended_time"` // Time spent with the freeze panel still up after observer freeze cam is done.
	SpecFreezeTimeLock                   int     `csgo:"spec_freeze_time_lock"`
	SpecFreezeDeathanimTime              int     `csgo:"spec_freeze_deathanim_time"`
	SvAccelerate                         float64 `csgo:"sv_accelerate"`                 // ( def. "10" ) client notify replicated
	SvStopspeed                          int     `csgo:"sv_stopspeed"`                  //
	SvAllowVotes                         int     `csgo:"sv_allow_votes"`                // Allow voting?
	SvAllowWaitCommand                   int     `csgo:"sv_allow_wait_command"`         // Allow or disallow the wait command on clients connected to this server.
	SvAlltalk                            int     `csgo:"sv_alltalk"`                    // Players can hear all other players' voice communication, no team restrictions
	SvAlternateticks                     int     `csgo:"sv_alternateticks"`             // If set, server only simulates entities on even numbered ticks.
	SvCheats                             int     `csgo:"sv_cheats"`                     // Allow cheats on server
	SvClockcorrectionMsecs               int     `csgo:"sv_clockcorrection_msecs"`      // The server tries to keep each player's m_nTickBase withing this many msecs of the server absolute tickcount
	SvConsistency                        int     `csgo:"sv_consistency"`                // Whether the server enforces file consistency for critical files
	SvContact                            int     `csgo:"sv_contact"`                    // Contact email for server sysop
	SvDamagePrintEnable                  int     `csgo:"sv_damage_print_enable"`        // Turn this off to disable the player's damage feed in the console after getting killed.
	SvDcFriendsReqd                      int     `csgo:"sv_dc_friends_reqd"`            // Set this to 0 to allow direct connects to a game in progress even if no presents
	SvDeadtalk                           int     `csgo:"sv_deadtalk"`                   // Dead players can speak (voice, text) to the living
	SvForcepreload                       int     `csgo:"sv_forcepreload"`               // Force server side preloading.
	SvFriction                           float64 `csgo:"sv_friction"`                   // World friction.
	SvFullAlltalk                        int     `csgo:"sv_full_alltalk"`               // Any player (including Spectator team) can speak to any other player
	SvGameinstructorDisable              int     `csgo:"sv_gameinstructor_disable"`     // Force all clients to disable their game instructors.
	SvIgnoregrenaderadio                 int     `csgo:"sv_ignoregrenaderadio"`         // Turn off Fire in the hole messages
	SvKickPlayersWithCooldown            int     `csgo:"sv_kick_players_with_cooldown"` // (0: do not kick; 1: kick Untrusted players; 2: kick players with any cooldown)
	SvKickBanDuration                    int     `csgo:"sv_kick_ban_duration"`          // How long should a kick ban from the server should last (in minutes)
	SvLan                                int     `csgo:"sv_lan"`                        // Server is a lan server ( no heartbeat, no authentication, no non-class C addresses )
	SvLogOnefile                         int     `csgo:"sv_log_onefile"`                // Log server information to only one file.
	SvLogbans                            int     `csgo:"sv_logbans"`                    // Log server bans in the server logs.
	SvLogecho                            int     `csgo:"sv_logecho"`                    // Echo log information to the console.
	SvLogfile                            int     `csgo:"sv_logfile"`                    // Log server information in the log file.
	SvLogflush                           int     `csgo:"sv_logflush"`                   // Flush the log file to disk on each write (slow).
	SvLogsdir                            string  `csgo:"sv_logsdir"`                    // Folder in the game directory where server logs will be stored.
	SvMaxrate                            int     `csgo:"sv_maxrate"`                    // min. 0.000000 max. 30000.000000 replicated  Max bandwidth rate allowed on server, 0 == unlimited
	SvMincmdrate                         int     `csgo:"sv_mincmdrate"`                 // This sets the minimum value for cl_cmdrate. 0 == unlimited.
	SvMinrate                            int     `csgo:"sv_minrate"`                    // Min bandwidth rate allowed on server, 0 == unlimited
	SvCompetitiveMinspec                 int     `csgo:"sv_competitive_minspec"`        // Enable to force certain client convars to minimum/maximum values to help prevent competitive advantages.
	SvCompetitiveOfficial5v5             int     `csgo:"sv_competitive_official_5v5"`   // Enable to force the server to show 5v5 scoreboards and allows spectators to see characters through walls.
	SvPausable                           int     `csgo:"sv_pausable"`                   // Is the server pausable.
	SvPure                               int     `csgo:"sv_pure"`
	SvPureKickClients                    int     `csgo:"sv_pure_kick_clients"`        // If set to 1, the server will kick clients with mismatching files. Otherwise, it will issue a warning to the client.
	SvPureTrace                          int     `csgo:"sv_pure_trace"`               // If set to 1, the server will print a message whenever a client is verifying a CR
	SvSpawnAfkBombDropTime               int     `csgo:"sv_spawn_afk_bomb_drop_time"` // Players that spawn and don't move for longer than sv_spawn_afk_bomb_drop_time (default 15 seconds) will automatically drop the bomb.
	SvSteamgroupExclusive                int     `csgo:"sv_steamgroup_exclusive"`     // If set, only members of Steam group will be able to join the server when it's empty, public people will be able to join the server only if it has players.
	SvVoiceenable                        int     `csgo:"sv_voiceenable"`
	SvAutoFullAlltalkDuringWarmupHalfEnd int     `csgo:"sv_auto_full_alltalk_during_warmup_half_end"`
}

// GetDefaultSettings Get Default Server Config
func GetDefaultSettings() CsgoServerSettings {
	return CsgoServerSettings{

		AmmoGrenadeLimit:                     1,
		AmmoGrenadeLimitFlashbang:            2,
		AmmoGrenadeLimitTotal:                4,
		BotQuota:                             0,
		CashPlayerBombDefused:                300,
		CashPlayerBombPlanted:                300,
		CashPlayerDamageHostage:              -30,
		CashPlayerInteractWithHostage:        150,
		CashPlayerKilledEnemyDefault:         300,
		CashPlayerKilledEnemyFactor:          1,
		CashPlayerKilledHostage:              -1000,
		CashPlayerKilledTeammate:             -300,
		CashPlayerRescuedHostage:             1000,
		CashTeamEliminationBombMap:           3250,
		CashTeamHostageAlive:                 150,
		CashTeamHostageInteraction:           150,
		CashTeamLoserBonus:                   1400,
		CashTeamLoserBonusConsecutiveRounds:  500,
		CashTeamPlantedBombButDefused:        800,
		CashTeamRescuedHostage:               750,
		CashTeamTerroristWinBomb:             3500,
		CashTeamWinByDefusingBomb:            3500,
		CashTeamWinByHostageRescue:           3500,
		CashPlayerGetKilled:                  0,
		CashPlayerRespawnAmount:              0,
		CashTeamEliminationHostageMapCt:      2000,
		CashTeamEliminationHostageMapT:       1000,
		CashTeamWinByTimeRunningOutBomb:      3250,
		CashTeamWinByTimeRunningOutHostage:   3250,
		FfDamageReductionGrenade:             0.85,
		FfDamageReductionBullets:             0.33,
		FfDamageReductionOther:               0.4,
		FfDamageReductionGrenadeSelf:         1,
		MpAfterroundmoney:                    0,
		MpAutokick:                           0,
		MpAutoteambalance:                    0,
		MpBuytime:                            15,
		MpC4timer:                            40,
		MpDeathDropDefuser:                   1,
		MpDeathDropGrenade:                   2,
		MpDeathDropGun:                       1,
		MpDefuserAllocation:                  0,
		MpDoWarmupPeriod:                     1,
		MpForcecamera:                        1,
		MpForcePickTime:                      160,
		MpFreeArmor:                          0,
		MpFreezetime:                         12,
		MpFriendlyfire:                       1,
		MpHalftime:                           1,
		MpHalftimeDuration:                   15,
		MpJoinGraceTime:                      30,
		MpLimitteams:                         0,
		MpLogdetail:                          3,
		MpMatchCanClinch:                     1,
		MpMatchEndRestart:                    1,
		MpMaxmoney:                           16000,
		MpMaxrounds:                          30,
		MpMolotovusedelay:                    0,
		MpPlayercashawards:                   1,
		MpPlayerid:                           0,
		MpPlayeridDelay:                      0.5,
		MpPlayeridHold:                       0.25,
		MpRoundRestartDelay:                  5,
		MpRoundtime:                          1.92,
		MpRoundtimeDefuse:                    1.92,
		MpSolidTeammates:                     1,
		MpStartmoney:                         800,
		MpTeamcashawards:                     1,
		MpTimelimit:                          0,
		MpTkpunish:                           0,
		MpWarmuptime:                         1,
		MpWeaponsAllowMapPlaced:              1,
		MpWeaponsAllowZeus:                   1,
		MpWinPanelDisplayTime:                15,
		MpOvertimeEnable:                     1,
		MpOvertimeMaxrounds:                  6,
		MpOvertimeStartmoney:                 10000,
		RconPassword:                         "NUEL",
		TvEnable:                             1,
		TvDelay:                              90,
		SpecFreezeTime:                       2.0,
		SpecFreezePanelExtendedTime:          0,
		SpecFreezeTimeLock:                   2,
		SpecFreezeDeathanimTime:              0,
		SvAccelerate:                         5.5,
		SvStopspeed:                          80,
		SvAllowVotes:                         0,
		SvAllowWaitCommand:                   0,
		SvAlltalk:                            0,
		SvAlternateticks:                     0,
		SvCheats:                             0,
		SvClockcorrectionMsecs:               15,
		SvConsistency:                        0,
		SvContact:                            0,
		SvDamagePrintEnable:                  0,
		SvDcFriendsReqd:                      0,
		SvDeadtalk:                           0,
		SvForcepreload:                       0,
		SvFriction:                           5.2,
		SvFullAlltalk:                        0,
		SvGameinstructorDisable:              1,
		SvIgnoregrenaderadio:                 0,
		SvKickPlayersWithCooldown:            0,
		SvKickBanDuration:                    0,
		SvLan:                                0,
		SvLogOnefile:                         0,
		SvLogbans:                            1,
		SvLogecho:                            1,
		SvLogfile:                            1,
		SvLogflush:                           0,
		SvLogsdir:                            "logfiles",
		SvMaxrate:                            0,
		SvMincmdrate:                         30,
		SvMinrate:                            20000,
		SvCompetitiveMinspec:                 1,
		SvCompetitiveOfficial5v5:             1,
		SvPausable:                           1,
		SvPure:                               1,
		SvPureKickClients:                    1,
		SvPureTrace:                          0,
		SvSpawnAfkBombDropTime:               30,
		SvSteamgroupExclusive:                0,
		SvVoiceenable:                        1,
		SvAutoFullAlltalkDuringWarmupHalfEnd: 0,
	}
}

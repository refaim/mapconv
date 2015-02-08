package homm2

type Serializer interface {
    UInt8(*uint8)
    UInt16(*uint16)
    UInt32(*uint32)
    CString(*string)
    IsReader() bool
}

type Color uint8
type ObjectType uint8
type CoodResourceType uint8
type VictoryCondition uint8
type LossCondition uint8
type Race uint8

type Tile struct {
    TileIndex     uint16
    ObjectName1   uint8
    IndexName1    uint8
    Quantity1     uint8
    Quantity2     uint8
    ObjectName2   uint8
    IndexName2    uint8
    Shape         uint8
    GeneralObject uint8
    IndexAddon    uint16
    UniqNumber1   uint32
    UniqNumber2   uint32
}

type TileAddon struct {
    IndexAddon   uint16
    ObjectNameN1 uint8
    IndexNameN1  uint8
    QuantityN    uint8
    ObjectNameN2 uint8
    IndexNameN2  uint8
    UniqNumberN1 uint32
    UniqNumberN2 uint32
}

type CoodCastle struct {
    X    uint8
    Y    uint8
    Type uint8
}

type CoodResource struct {
    X    uint8
    Y    uint8
    Type CoodResourceType
}

type Object struct {
    Size uint16
    Data []uint8
}

type Map struct {
    Magic                         uint32
    Complexity                    uint8
    Unknown1                      uint8
    Width                         uint8
    Height                        uint8
    PlayerPresence                [6]uint8
    PlayerHumanPlayable           [6]uint8
    PlayerComputerPlayable        [6]uint8
    PlayerCount                   uint8
    HumanOnlyPlayers              uint8
    NonComputerOnlyPlayers        uint8
    SpecialVictoryConditionsType  VictoryCondition
    ComputerCanWinViaSVC          uint8
    AllowNormalVictory            uint8
    SpecialVictoryConditionsData1 uint16
    SpecialLossConditionsType     LossCondition
    SpecialLossConditionsData1    uint16
    StartWithHero                 uint8
    PlayerRace                    [6]Race
    SpecialVictoryConditionsData2 uint16
    SpecialLossConditionsData2    uint16
    Unknown2                      [8]uint8
    Unknown3                      uint16
    Name                          [60]uint8
    Description                   [300]uint8
    RumorCount                    uint8
    EventCount                    uint8
    Unknown4                      [8]uint8
    Tiles                         []*Tile
    TileAddonCount                uint32
    TileAddons                    []*TileAddon
    CoodCastles                   [72]CoodCastle
    CoodResourceKingdoms          [144]CoodResource
    ObeliskCount                  uint8
    RumorObjectIds                []uint16
    EventObjectIds                []uint16
    ObjectCount                   uint16
    Objects                       []*Object
    Unknown5                      uint32
}

type Castle struct {
    Color              Color
    UseCustomBuildings uint8
    Buildings          uint16
    Dwellings          uint16
    MagicTower         uint8
    UseCustomTroops    uint8
    MonsterType        [5]uint8
    MonsterCount       [5]uint16
    HasCapitan         uint8
    UseCustomName      uint8
    CastleName         [13]uint8
    Type               uint8
    IsCastle           uint8
    ForbidCastle       uint8
    Unknown1           [29]uint8
}

type Hero struct {
    Unknown1       uint8
    CustomTroops   uint8
    MonsterType    [5]uint8
    MonsterCount   [5]uint16
    CustomPortrate uint8
    Portrate       uint8
    Artifacts      [3]uint8
    Unknown2       uint8
    Experience     uint32
    CustomSkills   uint8
    SkillTypes     [8]uint8
    SkillLevels    [8]uint8
    Unknown3       [15]uint8
    CustomName     uint8
    Name           [13]uint8
    Patrol         uint8
    CountSquare    uint8
}

type Info struct {
    Id       uint8
    Unknown1 [8]uint8
    Text     string
}

type EventCoord struct {
    Id       uint8
    Wood     uint32
    Mercury  uint32
    Ore      uint32
    Sulfur   uint32
    Crystal  uint32
    Gems     uint32
    Golds    uint32
    Artifact uint16
    Computer uint8
    Cancel   uint8
    Unknown1 [8]uint8
    Blue     uint8
    Green    uint8
    Red      uint8
    Yellow   uint8
    Orange   uint8
    Purple   uint8
    Text     string
}

type EventDay struct {
    Id         uint8
    Wood       uint32
    Mercury    uint32
    Ore        uint32
    Sulfur     uint32
    Crystal    uint32
    Gems       uint32
    Golds      uint32
    Artifact   uint16
    Computer   uint16
    First      uint16
    Subsequent uint16
    Unknown1   [6]uint8
    Zero       uint8
    Blue       uint8
    Green      uint8
    Red        uint8
    Yellow     uint8
    Orange     uint8
    Purple     uint8
    Text       string
}

type Rumor struct {
    Id       uint8
    Unknown1 [7]uint8
    Text     string
}

type Riddle struct {
    Id       uint8
    Wood     uint32
    Mercury  uint32
    Ore      uint32
    Sulfur   uint32
    Crystal  uint32
    Gems     uint32
    Golds    uint32
    Artifact uint16
    Count    uint8
    Answer1  [13]uint8
    Answer2  [13]uint8
    Answer3  [13]uint8
    Answer4  [13]uint8
    Answer5  [13]uint8
    Answer6  [13]uint8
    Answer7  [13]uint8
    Answer8  [13]uint8
    Text     string
}

const (
    ColorBlue    Color = 0
    ColorGreen         = 1
    ColorRed           = 2
    ColorYellow        = 3
    ColorOrange        = 4
    ColorPurple        = 5
    ColorUnknown       = 255
)

const (
    BuildingsThiefGuild    = (1 << 1)
    BuildingsTavern        = (1 << 2)
    BuildingsShipyard      = (1 << 3)
    BuildingsWell          = (1 << 4)
    BuildingsStatue        = (1 << 6)
    BuildingsLeftTurret    = (1 << 7)
    BuildingsRightTurret   = (1 << 8)
    BuildingsMarketplace   = (1 << 9)
    BuildingsFarm          = (1 << 10)
    BuildingsMoat          = (1 << 11)
    BuildingsFortification = (1 << 12)
)

const (
    DwellingsLevel1         = (1 << 3)
    DwellingsLevel2         = (1 << 4)
    DwellingsLevel3         = (1 << 5)
    DwellingsLevel4         = (1 << 6)
    DwellingsLevel5         = (1 << 7)
    DwellingsLevel6         = (1 << 8)
    DwellingsUpgradedLevel2 = (1 << 9)
    DwellingsUpgradedLevel3 = (1 << 10)
    DwellingsUpgradedLevel4 = (1 << 11)
    DwellingsUpgradedLevel5 = (1 << 12)
    DwellingsUpgradedLevel6 = (1 << 13)
)

const (
    ObjectTypeNone                   ObjectType = 0
    ObjectTypeAlchemyLab                        = 129
    ObjectTypeSkeleton                          = 132
    ObjectTypeDaemonCave                        = 133
    ObjectTypeFaerieRing                        = 135
    ObjectTypeGazebo                            = 138
    ObjectTypeGraveyard                         = 140
    ObjectTypeArcherHouse                       = 141
    ObjectTypeDwarfCott                         = 143
    ObjectTypePeasantHut                        = 144
    ObjectTypeDragonCity                        = 148
    ObjectTypeLighthouse                        = 149
    ObjectTypeWaterWheel                        = 22
    ObjectTypeMines                             = 151
    ObjectTypeObelisk                           = 153
    ObjectTypeOasis                             = 154
    ObjectTypeCoast                             = 28
    ObjectTypeSawmill                           = 157
    ObjectTypeOracle                            = 158
    ObjectTypeShipwreck                         = 160
    ObjectTypeDesertTent                        = 162
    ObjectTypeCastle                            = 163
    ObjectTypeStoneLights                       = 164
    ObjectTypeWagonCamp                         = 37
    ObjectTypeWindmill                          = 168
    ObjectTypeRandomTown                        = 176
    ObjectTypeRandomCastle                      = 177
    ObjectTypeShrub2                            = 56
    ObjectTypeNothingSpecial                    = 57
    ObjectTypeWatchTower                        = 186
    ObjectTypeTreeHouse                         = 187
    ObjectTypeTreeCity                          = 188
    ObjectTypeRuins                             = 189
    ObjectTypeFort                              = 190
    ObjectTypeTradingPost                       = 191
    ObjectTypeAbandonedMine                     = 192
    ObjectTypeTreeKnowledge                     = 196
    ObjectTypeDoctorHut                         = 69
    ObjectTypeTemple                            = 198
    ObjectTypeHillFort                          = 199
    ObjectTypeHalflingHole                      = 200
    ObjectTypeMercenaryCamp                     = 201
    ObjectTypePyramid                           = 204
    ObjectTypeCitydead                          = 205
    ObjectTypeExcavation                        = 206
    ObjectTypeSphinx                            = 207
    ObjectTypeTarPit                            = 81
    ObjectTypeArtesianSpring                    = 210
    ObjectTypeTrollBridge                       = 211
    ObjectTypeWateringHole                      = 212
    ObjectTypeWitchShut                         = 213
    ObjectTypeXanadu                            = 214
    ObjectTypeCave                              = 215
    ObjectTypeMagellanMaps                      = 217
    ObjectTypeDerelictShip                      = 219
    ObjectTypeMagicWell                         = 222
    ObjectTypeObservationTower                  = 224
    ObjectTypeFreemanFoundry                    = 225
    ObjectTypeTrees                             = 99
    ObjectTypeMounts                            = 100
    ObjectTypeVolcano                           = 101
    ObjectTypeFlowers                           = 102
    ObjectTypeStones                            = 103
    ObjectTypeWaterLake                         = 104
    ObjectTypeMandrake                          = 105
    ObjectTypeDeadTree                          = 106
    ObjectTypeStump                             = 107
    ObjectTypeCrater                            = 108
    ObjectTypeCactus                            = 109
    ObjectTypeMound                             = 110
    ObjectTypeDune                              = 111
    ObjectTypeLavaPool                          = 112
    ObjectTypeShrub                             = 113
    ObjectTypeArena                             = 242
    ObjectTypeBarrowMounds                      = 243
    ObjectTypeMermaid                           = 236
    ObjectTypeSirens                            = 237
    ObjectTypeHutMagi                           = 118
    ObjectTypeEyeMagi                           = 119
    ObjectTypeTravellerTent                     = 248
    ObjectTypeJail                              = 251
    ObjectTypeFireAltar                         = 252
    ObjectTypeAirAltar                          = 253
    ObjectTypeEarthAltar                        = 254
    ObjectTypeWaterAltar                        = 255
    ObjectTypeWaterChest                        = 128
    ObjectTypeSign                              = 130
    ObjectTypeBuoy                              = 131
    ObjectTypeTreasureChest                     = 134
    ObjectTypeCampfire                          = 136
    ObjectTypeFountain                          = 137
    ObjectTypeAncientLamp                       = 139
    ObjectTypeGoblinHut                         = 142
    ObjectTypeEvent                             = 147
    ObjectTypeWaterwheel                        = 150
    ObjectTypeMonster                           = 152
    ObjectTypeResource                          = 155
    ObjectTypeShrine1                           = 159
    ObjectTypeWagoncamp                         = 165
    ObjectTypeWhirlpool                         = 167
    ObjectTypeArtifact                          = 169
    ObjectTypeBoat                              = 171
    ObjectTypeRandomUltimateArtifact            = 172
    ObjectTypeRandomArtifact                    = 173
    ObjectTypeRandomResource                    = 174
    ObjectTypeRandomMonster                     = 175
    ObjectTypeRandomMonster1                    = 179
    ObjectTypeRandomMonster2                    = 180
    ObjectTypeRandomMonster3                    = 181
    ObjectTypeRandomMonster4                    = 182
    ObjectTypeHeroes                            = 183
    ObjectTypeThatchedHut                       = 193
    ObjectTypeStandingStones                    = 194
    ObjectTypeIdol                              = 195
    ObjectTypeDoctorhut                         = 197
    ObjectTypeShrine2                           = 202
    ObjectTypeShrine3                           = 203
    ObjectTypeWagon                             = 208
    ObjectTypeLeanto                            = 216
    ObjectTypeFlotsam                           = 218
    ObjectTypeShipwreckSurviror                 = 220
    ObjectTypeBottle                            = 221
    ObjectTypeMagicGarden                       = 223
    ObjectTypeReefs                             = 233
    ObjectTypeAlchemyTower                      = 240
    ObjectTypeStables                           = 241
    ObjectTypeHutmagi                           = 238
    ObjectTypeEyemagi                           = 239
    ObjectTypeRandomArtifact1                   = 244
    ObjectTypeRandomArtifact2                   = 245
    ObjectTypeRandomArtifact3                   = 246
    ObjectTypeBarrier                           = 247
)

const (
    CoodCastleTypeKnight      = 0
    CoodCastleTypeBarbarian   = 1
    CoodCastleTypeSorceress   = 2
    CoodCastleTypeWarlock     = 3
    CoodCastleTypeWizard      = 4
    CoodCastleTypeNecromancer = 5
    CoodCastleTypeRandom      = 6
    CoodCastleTypeCastle      = 128
)

const (
    CoodResourceTypeWoodMine      CoodResourceType = 0
    CoodResourceTypeMercuryMin                     = 1
    CoodResourceTypeOreMine                        = 2
    CoodResourceTypeSulfurMine                     = 3
    CoodResourceTypeCrystalMine                    = 4
    CoodResourceTypeGemsMine                       = 5
    CoodResourceTypeGoldMine                       = 6
    CoodResourceTypeLightHouse                     = 100
    CoodResourceTypeDragonCity                     = 101
    CoodResourceTypeAbandonedMine                  = 103
)

const (
    VictoryConditionAll      VictoryCondition = 0
    VictoryConditionTown                      = 1
    VictoryConditionHero                      = 2
    VictoryConditionArtifact                  = 3
    VictoryConditionSide                      = 4
    VictoryConditionGold                      = 5
)

const (
    LossConditionAll  LossCondition = 0
    LossConditionTown               = 1
    LossConditionHero               = 2
    LossConditionTime               = 3
)

const (
    RaceKnight      Race = 0
    RaceBarbarian        = 1
    RaceSorceress        = 2
    RaceWarlock          = 3
    RaceWizard           = 4
    RaceNecromancer      = 5
    RaceMultiple         = 6
    RaceRandom           = 7
    RaceNone             = 255
)

func (s *Tile) Serialize(b Serializer) {
    b.UInt16((&(s.TileIndex)))
    b.UInt8((&(s.ObjectName1)))
    b.UInt8((&(s.IndexName1)))
    b.UInt8((&(s.Quantity1)))
    b.UInt8((&(s.Quantity2)))
    b.UInt8((&(s.ObjectName2)))
    b.UInt8((&(s.IndexName2)))
    b.UInt8((&(s.Shape)))
    b.UInt8((&(s.GeneralObject)))
    b.UInt16((&(s.IndexAddon)))
    b.UInt32((&(s.UniqNumber1)))
    b.UInt32((&(s.UniqNumber2)))
}

func (s *TileAddon) Serialize(b Serializer) {
    b.UInt16((&(s.IndexAddon)))
    b.UInt8((&(s.ObjectNameN1)))
    b.UInt8((&(s.IndexNameN1)))
    b.UInt8((&(s.QuantityN)))
    b.UInt8((&(s.ObjectNameN2)))
    b.UInt8((&(s.IndexNameN2)))
    b.UInt32((&(s.UniqNumberN1)))
    b.UInt32((&(s.UniqNumberN2)))
}

func (s *CoodCastle) Serialize(b Serializer) {
    b.UInt8((&(s.X)))
    b.UInt8((&(s.Y)))
    b.UInt8((&(s.Type)))
}

func (s *CoodResource) Serialize(b Serializer) {
    b.UInt8((&(s.X)))
    b.UInt8((&(s.Y)))
    b.UInt8((*uint8)((&(s.Type))))
}

func (s *Object) Serialize(b Serializer) {
    b.UInt16((&(s.Size)))
    size1 := int(s.Size)
    if b.IsReader() {        *(&(s.Data)) = make([]uint8, size1)
    }
    for i1 := (0); i1 < size1; i1++ {
        b.UInt8((&(*(&(s.Data)))[i1]))
    }
}

func (s *Map) Serialize(b Serializer) {
    b.UInt32((&(s.Magic)))
    b.UInt8((&(s.Complexity)))
    b.UInt8((&(s.Unknown1)))
    b.UInt8((&(s.Width)))
    b.UInt8((&(s.Height)))
    size2 := 6
    for i2 := (0); i2 < size2; i2++ {
        b.UInt8((&(*(&(s.PlayerPresence)))[i2]))
    }
    size3 := 6
    for i3 := (0); i3 < size3; i3++ {
        b.UInt8((&(*(&(s.PlayerHumanPlayable)))[i3]))
    }
    size4 := 6
    for i4 := (0); i4 < size4; i4++ {
        b.UInt8((&(*(&(s.PlayerComputerPlayable)))[i4]))
    }
    b.UInt8((&(s.PlayerCount)))
    b.UInt8((&(s.HumanOnlyPlayers)))
    b.UInt8((&(s.NonComputerOnlyPlayers)))
    b.UInt8((*uint8)((&(s.SpecialVictoryConditionsType))))
    b.UInt8((&(s.ComputerCanWinViaSVC)))
    b.UInt8((&(s.AllowNormalVictory)))
    b.UInt16((&(s.SpecialVictoryConditionsData1)))
    b.UInt8((*uint8)((&(s.SpecialLossConditionsType))))
    b.UInt16((&(s.SpecialLossConditionsData1)))
    b.UInt8((&(s.StartWithHero)))
    size5 := 6
    for i5 := (0); i5 < size5; i5++ {
        b.UInt8((*uint8)((&(*(&(s.PlayerRace)))[i5])))
    }
    b.UInt16((&(s.SpecialVictoryConditionsData2)))
    b.UInt16((&(s.SpecialLossConditionsData2)))
    size6 := 8
    for i6 := (0); i6 < size6; i6++ {
        b.UInt8((&(*(&(s.Unknown2)))[i6]))
    }
    b.UInt16((&(s.Unknown3)))
    size7 := 60
    for i7 := (0); i7 < size7; i7++ {
        b.UInt8((&(*(&(s.Name)))[i7]))
    }
    size8 := 300
    for i8 := (0); i8 < size8; i8++ {
        b.UInt8((&(*(&(s.Description)))[i8]))
    }
    b.UInt8((&(s.RumorCount)))
    b.UInt8((&(s.EventCount)))
    size9 := 8
    for i9 := (0); i9 < size9; i9++ {
        b.UInt8((&(*(&(s.Unknown4)))[i9]))
    }
    size10 := (int(s.Width) * int(s.Height))
    if b.IsReader() {        *(&(s.Tiles)) = make([]*Tile, size10)
        for i10 := (0); i10 < size10; i10++ {
            *(&(*(&(s.Tiles)))[i10]) = &Tile{}
        }
    }
    for i10 := (0); i10 < size10; i10++ {
        (&(*(&(s.Tiles)))[i10]).Serialize(b)
    }
    b.UInt32((&(s.TileAddonCount)))
    size11 := int(s.TileAddonCount)
    if b.IsReader() {        *(&(s.TileAddons)) = make([]*TileAddon, size11)
        for i11 := (0); i11 < size11; i11++ {
            *(&(*(&(s.TileAddons)))[i11]) = &TileAddon{}
        }
    }
    for i11 := (0); i11 < size11; i11++ {
        (&(*(&(s.TileAddons)))[i11]).Serialize(b)
    }
    size12 := 72
    for i12 := (0); i12 < size12; i12++ {
        (&(*(&(s.CoodCastles)))[i12]).Serialize(b)
    }
    size13 := 144
    for i13 := (0); i13 < size13; i13++ {
        (&(*(&(s.CoodResourceKingdoms)))[i13]).Serialize(b)
    }
    b.UInt8((&(s.ObeliskCount)))
    size14 := int(s.RumorCount)
    if b.IsReader() {        *(&(s.RumorObjectIds)) = make([]uint16, size14)
    }
    for i14 := (0); i14 < size14; i14++ {
        b.UInt16((&(*(&(s.RumorObjectIds)))[i14]))
    }
    size15 := int(s.EventCount)
    if b.IsReader() {        *(&(s.EventObjectIds)) = make([]uint16, size15)
    }
    for i15 := (0); i15 < size15; i15++ {
        b.UInt16((&(*(&(s.EventObjectIds)))[i15]))
    }
    b.UInt16((&(s.ObjectCount)))
    size16 := int(s.ObjectCount)
    if b.IsReader() {        *(&(s.Objects)) = make([]*Object, size16)
        for i16 := (0); i16 < size16; i16++ {
            *(&(*(&(s.Objects)))[i16]) = &Object{}
        }
    }
    for i16 := (0); i16 < size16; i16++ {
        (&(*(&(s.Objects)))[i16]).Serialize(b)
    }
    b.UInt32((&(s.Unknown5)))
}

func (s *Castle) Serialize(b Serializer) {
    b.UInt8((*uint8)((&(s.Color))))
    b.UInt8((&(s.UseCustomBuildings)))
    b.UInt16((&(s.Buildings)))
    b.UInt16((&(s.Dwellings)))
    b.UInt8((&(s.MagicTower)))
    b.UInt8((&(s.UseCustomTroops)))
    size17 := 5
    for i17 := (0); i17 < size17; i17++ {
        b.UInt8((&(*(&(s.MonsterType)))[i17]))
    }
    size18 := 5
    for i18 := (0); i18 < size18; i18++ {
        b.UInt16((&(*(&(s.MonsterCount)))[i18]))
    }
    b.UInt8((&(s.HasCapitan)))
    b.UInt8((&(s.UseCustomName)))
    size19 := 13
    for i19 := (0); i19 < size19; i19++ {
        b.UInt8((&(*(&(s.CastleName)))[i19]))
    }
    b.UInt8((&(s.Type)))
    b.UInt8((&(s.IsCastle)))
    b.UInt8((&(s.ForbidCastle)))
    size20 := 29
    for i20 := (0); i20 < size20; i20++ {
        b.UInt8((&(*(&(s.Unknown1)))[i20]))
    }
}

func (s *Hero) Serialize(b Serializer) {
    b.UInt8((&(s.Unknown1)))
    b.UInt8((&(s.CustomTroops)))
    size21 := 5
    for i21 := (0); i21 < size21; i21++ {
        b.UInt8((&(*(&(s.MonsterType)))[i21]))
    }
    size22 := 5
    for i22 := (0); i22 < size22; i22++ {
        b.UInt16((&(*(&(s.MonsterCount)))[i22]))
    }
    b.UInt8((&(s.CustomPortrate)))
    b.UInt8((&(s.Portrate)))
    size23 := 3
    for i23 := (0); i23 < size23; i23++ {
        b.UInt8((&(*(&(s.Artifacts)))[i23]))
    }
    b.UInt8((&(s.Unknown2)))
    b.UInt32((&(s.Experience)))
    b.UInt8((&(s.CustomSkills)))
    size24 := 8
    for i24 := (0); i24 < size24; i24++ {
        b.UInt8((&(*(&(s.SkillTypes)))[i24]))
    }
    size25 := 8
    for i25 := (0); i25 < size25; i25++ {
        b.UInt8((&(*(&(s.SkillLevels)))[i25]))
    }
    size26 := 15
    for i26 := (0); i26 < size26; i26++ {
        b.UInt8((&(*(&(s.Unknown3)))[i26]))
    }
    b.UInt8((&(s.CustomName)))
    size27 := 13
    for i27 := (0); i27 < size27; i27++ {
        b.UInt8((&(*(&(s.Name)))[i27]))
    }
    b.UInt8((&(s.Patrol)))
    b.UInt8((&(s.CountSquare)))
}

func (s *Info) Serialize(b Serializer) {
    b.UInt8((&(s.Id)))
    size28 := 8
    for i28 := (0); i28 < size28; i28++ {
        b.UInt8((&(*(&(s.Unknown1)))[i28]))
    }
    b.CString((&(s.Text)))
}

func (s *EventCoord) Serialize(b Serializer) {
    b.UInt8((&(s.Id)))
    b.UInt32((&(s.Wood)))
    b.UInt32((&(s.Mercury)))
    b.UInt32((&(s.Ore)))
    b.UInt32((&(s.Sulfur)))
    b.UInt32((&(s.Crystal)))
    b.UInt32((&(s.Gems)))
    b.UInt32((&(s.Golds)))
    b.UInt16((&(s.Artifact)))
    b.UInt8((&(s.Computer)))
    b.UInt8((&(s.Cancel)))
    size29 := 8
    for i29 := (0); i29 < size29; i29++ {
        b.UInt8((&(*(&(s.Unknown1)))[i29]))
    }
    b.UInt8((&(s.Blue)))
    b.UInt8((&(s.Green)))
    b.UInt8((&(s.Red)))
    b.UInt8((&(s.Yellow)))
    b.UInt8((&(s.Orange)))
    b.UInt8((&(s.Purple)))
    b.CString((&(s.Text)))
}

func (s *EventDay) Serialize(b Serializer) {
    b.UInt8((&(s.Id)))
    b.UInt32((&(s.Wood)))
    b.UInt32((&(s.Mercury)))
    b.UInt32((&(s.Ore)))
    b.UInt32((&(s.Sulfur)))
    b.UInt32((&(s.Crystal)))
    b.UInt32((&(s.Gems)))
    b.UInt32((&(s.Golds)))
    b.UInt16((&(s.Artifact)))
    b.UInt16((&(s.Computer)))
    b.UInt16((&(s.First)))
    b.UInt16((&(s.Subsequent)))
    size30 := 6
    for i30 := (0); i30 < size30; i30++ {
        b.UInt8((&(*(&(s.Unknown1)))[i30]))
    }
    b.UInt8((&(s.Zero)))
    b.UInt8((&(s.Blue)))
    b.UInt8((&(s.Green)))
    b.UInt8((&(s.Red)))
    b.UInt8((&(s.Yellow)))
    b.UInt8((&(s.Orange)))
    b.UInt8((&(s.Purple)))
    b.CString((&(s.Text)))
}

func (s *Rumor) Serialize(b Serializer) {
    b.UInt8((&(s.Id)))
    size31 := 7
    for i31 := (0); i31 < size31; i31++ {
        b.UInt8((&(*(&(s.Unknown1)))[i31]))
    }
    b.CString((&(s.Text)))
}

func (s *Riddle) Serialize(b Serializer) {
    b.UInt8((&(s.Id)))
    b.UInt32((&(s.Wood)))
    b.UInt32((&(s.Mercury)))
    b.UInt32((&(s.Ore)))
    b.UInt32((&(s.Sulfur)))
    b.UInt32((&(s.Crystal)))
    b.UInt32((&(s.Gems)))
    b.UInt32((&(s.Golds)))
    b.UInt16((&(s.Artifact)))
    b.UInt8((&(s.Count)))
    size32 := 13
    for i32 := (0); i32 < size32; i32++ {
        b.UInt8((&(*(&(s.Answer1)))[i32]))
    }
    size33 := 13
    for i33 := (0); i33 < size33; i33++ {
        b.UInt8((&(*(&(s.Answer2)))[i33]))
    }
    size34 := 13
    for i34 := (0); i34 < size34; i34++ {
        b.UInt8((&(*(&(s.Answer3)))[i34]))
    }
    size35 := 13
    for i35 := (0); i35 < size35; i35++ {
        b.UInt8((&(*(&(s.Answer4)))[i35]))
    }
    size36 := 13
    for i36 := (0); i36 < size36; i36++ {
        b.UInt8((&(*(&(s.Answer5)))[i36]))
    }
    size37 := 13
    for i37 := (0); i37 < size37; i37++ {
        b.UInt8((&(*(&(s.Answer6)))[i37]))
    }
    size38 := 13
    for i38 := (0); i38 < size38; i38++ {
        b.UInt8((&(*(&(s.Answer7)))[i38]))
    }
    size39 := 13
    for i39 := (0); i39 < size39; i39++ {
        b.UInt8((&(*(&(s.Answer8)))[i39]))
    }
    b.CString((&(s.Text)))
}


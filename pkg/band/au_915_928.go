// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package band

import (
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

//revive:disable:var-naming

var au_915_928 Band

// AU_915_928 is the ID of the Australian band
const AU_915_928 = "AU_915_928"

//revive:enable:var-naming

func init() {
	uplinkChannels := make([]Channel, 0, 72)
	for i := 0; i < 64; i++ {
		uplinkChannels = append(uplinkChannels, Channel{
			Frequency:   uint64(915200000 + 200000*i),
			MinDataRate: 0,
			MaxDataRate: 3,
		})
	}
	for i := 0; i < 8; i++ {
		uplinkChannels = append(uplinkChannels, Channel{
			Frequency:   uint64(915900000 + 1600000*i),
			MinDataRate: 4,
			MaxDataRate: 4,
		})
	}

	downlinkChannels := make([]Channel, 0, 8)
	for i := 0; i < 8; i++ {
		downlinkChannels = append(downlinkChannels, Channel{
			Frequency:   uint64(923300000 + 600000*i),
			MinDataRate: 8,
			MaxDataRate: 13,
		})
	}

	downlinkDRTable := [7][6]ttnpb.DataRateIndex{
		{8, 8, 8, 8, 8, 8},
		{9, 8, 8, 8, 8, 8},
		{10, 9, 8, 8, 8, 8},
		{11, 10, 9, 8, 8, 8},
		{12, 11, 10, 9, 8, 8},
		{13, 12, 11, 10, 9, 8},
		{13, 13, 12, 11, 10, 9},
	}

	au_915_928 = Band{
		ID: AU_915_928,

		MaxUplinkChannels: 72,
		UplinkChannels:    uplinkChannels,

		MaxDownlinkChannels: 8,
		DownlinkChannels:    downlinkChannels,

		// See Radiocommunications (Low Interference Potential Devices) Class Licence 2015
		SubBands: []SubBandParameters{
			{
				MinFrequency: 915000000,
				MaxFrequency: 928000000,
				DutyCycle:    1,
				MaxEIRP:      30,
			},
		},

		DataRates: map[ttnpb.DataRateIndex]DataRate{
			0: makeLoRaDataRate(12, 125000, makeConstMaxMACPayloadSizeFunc(59)),
			1: makeLoRaDataRate(11, 125000, makeConstMaxMACPayloadSizeFunc(59)),
			2: makeLoRaDataRate(10, 125000, makeConstMaxMACPayloadSizeFunc(59)),
			3: makeLoRaDataRate(9, 125000, makeConstMaxMACPayloadSizeFunc(123)),
			4: makeLoRaDataRate(8, 125000, makeConstMaxMACPayloadSizeFunc(230)),
			5: makeLoRaDataRate(7, 125000, makeConstMaxMACPayloadSizeFunc(230)),
			6: makeLoRaDataRate(8, 500000, makeConstMaxMACPayloadSizeFunc(230)),

			8:  makeLoRaDataRate(12, 500000, makeConstMaxMACPayloadSizeFunc(41)),
			9:  makeLoRaDataRate(11, 500000, makeConstMaxMACPayloadSizeFunc(117)),
			10: makeLoRaDataRate(10, 500000, makeConstMaxMACPayloadSizeFunc(230)),
			11: makeLoRaDataRate(9, 500000, makeConstMaxMACPayloadSizeFunc(230)),
			12: makeLoRaDataRate(8, 500000, makeConstMaxMACPayloadSizeFunc(230)),
			13: makeLoRaDataRate(7, 500000, makeConstMaxMACPayloadSizeFunc(230)),
		},
		MaxADRDataRateIndex: 5,

		ReceiveDelay1:    defaultReceiveDelay1,
		ReceiveDelay2:    defaultReceiveDelay2,
		JoinAcceptDelay1: defaultJoinAcceptDelay1,
		JoinAcceptDelay2: defaultJoinAcceptDelay2,
		MaxFCntGap:       defaultMaxFCntGap,
		ADRAckLimit:      defaultADRAckLimit,
		ADRAckDelay:      defaultADRAckDelay,
		MinAckTimeout:    defaultAckTimeout - defaultAckTimeoutMargin,
		MaxAckTimeout:    defaultAckTimeout + defaultAckTimeoutMargin,

		DefaultMaxEIRP: 30,
		TxOffset: func() [16]float32 {
			offset := [16]float32{}
			for i := 0; i < 15; i++ {
				offset[i] = float32(0 - 2*i)
			}
			return offset
		}(),
		MaxTxPowerIndex: 14,

		LoRaCodingRate: "4/5",

		FreqMultiplier:   100,
		ImplementsCFList: true,
		CFListType:       ttnpb.CFListType_CHANNEL_MASKS,

		Rx1Channel: channelIndexModulo(8),
		Rx1DataRate: func(idx ttnpb.DataRateIndex, offset uint32, _ bool) (ttnpb.DataRateIndex, error) {
			if idx > 6 {
				return 0, errDataRateIndexTooHigh.WithAttributes("max", 6)
			}
			if offset > 5 {
				return 0, errDataRateOffsetTooHigh.WithAttributes("max", 5)
			}
			return downlinkDRTable[idx][offset], nil
		},

		GenerateChMasks: makeGenerateChMask72(true),
		ParseChMask:     parseChMask72,

		DefaultRx2Parameters: Rx2Parameters{8, 923300000},

		Beacon: Beacon{
			DataRateIndex:    8,
			CodingRate:       "4/5",
			ComputeFrequency: makeBeaconFrequencyFunc(usAuBeaconFrequencies),
		},

		TxParamSetupReqSupport: true,

		// No LoRaWAN Regional Parameters 1.0
		regionalParameters1_0_1:     bandIdentity,
		regionalParameters1_0_2RevA: auDataRates1_0_2,
		regionalParameters1_0_2RevB: composeSwaps(
			disableChMaskCntl51_0_2,
			disableTxParamSetupReq,
			makeSetMaxTxPowerIndexFunc(10),
		),
		regionalParameters1_0_3RevA: makeSetMaxTxPowerIndexFunc(15),
		regionalParameters1_1RevA:   bandIdentity,
	}
	All[AU_915_928] = au_915_928
}

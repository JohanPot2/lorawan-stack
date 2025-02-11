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

var us_902_928 Band

// US_902_928 is the ID of the US frequency plan
const US_902_928 = "US_902_928"

//revive:enable:var-naming

func init() {
	uplinkChannels := make([]Channel, 0, 72)
	for i := 0; i < 64; i++ {
		uplinkChannels = append(uplinkChannels, Channel{
			Frequency:   uint64(902300000 + 200000*i),
			MinDataRate: 0,
			MaxDataRate: 3,
		})
	}
	for i := 0; i < 8; i++ {
		uplinkChannels = append(uplinkChannels, Channel{
			Frequency:   uint64(903000000 + 1600000*i),
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

	downlinkDRTable := [5][4]ttnpb.DataRateIndex{
		{10, 9, 8, 8},
		{11, 10, 9, 8},
		{12, 11, 10, 9},
		{13, 12, 11, 10},
		{13, 13, 12, 11},
	}

	us_902_928 = Band{
		ID: US_902_928,

		MaxUplinkChannels: 72,
		UplinkChannels:    uplinkChannels,

		MaxDownlinkChannels: 8,
		DownlinkChannels:    downlinkChannels,

		// As per FCC Rules for Unlicensed Wireless Equipment operating in the ISM bands
		SubBands: []SubBandParameters{
			{
				MinFrequency: 902300000,
				MaxFrequency: 914900000,
				DutyCycle:    1,
				MaxEIRP:      21.0 + eirpDelta,
			},
			{
				MinFrequency: 923300000,
				MaxFrequency: 927500000,
				DutyCycle:    1,
				MaxEIRP:      26.0 + eirpDelta,
			},
		},

		DataRates: map[ttnpb.DataRateIndex]DataRate{
			0: makeLoRaDataRate(10, 125000, makeConstMaxMACPayloadSizeFunc(19)),
			1: makeLoRaDataRate(9, 125000, makeConstMaxMACPayloadSizeFunc(61)),
			2: makeLoRaDataRate(8, 125000, makeConstMaxMACPayloadSizeFunc(133)),
			3: makeLoRaDataRate(7, 125000, makeConstMaxMACPayloadSizeFunc(250)),
			4: makeLoRaDataRate(8, 500000, makeConstMaxMACPayloadSizeFunc(250)),

			8:  makeLoRaDataRate(12, 500000, makeConstMaxMACPayloadSizeFunc(41)),
			9:  makeLoRaDataRate(11, 500000, makeConstMaxMACPayloadSizeFunc(117)),
			10: makeLoRaDataRate(10, 500000, makeConstMaxMACPayloadSizeFunc(230)),
			11: makeLoRaDataRate(9, 500000, makeConstMaxMACPayloadSizeFunc(230)),
			12: makeLoRaDataRate(8, 500000, makeConstMaxMACPayloadSizeFunc(230)),
			13: makeLoRaDataRate(7, 500000, makeConstMaxMACPayloadSizeFunc(230)),
		},
		MaxADRDataRateIndex: 3,

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

		Rx1Channel: channelIndexModulo(8),
		Rx1DataRate: func(idx ttnpb.DataRateIndex, offset uint32, _ bool) (ttnpb.DataRateIndex, error) {
			if idx > 4 {
				return 0, errDataRateIndexTooHigh.WithAttributes("max", 4)
			}
			if offset > 3 {
				return 0, errDataRateOffsetTooHigh.WithAttributes("max", 3)
			}
			return downlinkDRTable[idx][offset], nil
		},

		GenerateChMasks: makeGenerateChMask72(true),
		ParseChMask:     parseChMask72,

		LoRaCodingRate: "4/5",

		FreqMultiplier:   100,
		ImplementsCFList: true,
		CFListType:       ttnpb.CFListType_CHANNEL_MASKS,

		DefaultRx2Parameters: Rx2Parameters{8, 923300000},

		Beacon: Beacon{
			DataRateIndex:    8,
			CodingRate:       "4/5",
			ComputeFrequency: makeBeaconFrequencyFunc(usAuBeaconFrequencies),
		},

		regionalParameters1_0:       bandIdentity,
		regionalParameters1_0_1:     bandIdentity,
		regionalParameters1_0_2RevA: usBeacon1_0_2,
		regionalParameters1_0_2RevB: composeSwaps(
			disableCFList1_0_2,
			disableChMaskCntl51_0_2,
			makeSetMaxTxPowerIndexFunc(10),
		),
		regionalParameters1_0_3RevA: makeSetMaxTxPowerIndexFunc(15),
		regionalParameters1_1RevA:   bandIdentity,
	}
	All[US_902_928] = us_902_928
}

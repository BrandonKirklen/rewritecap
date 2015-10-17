// Copyright 2014-2015 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os"
	"time"
)

//
// -----------------------------------------------------------------------------
// changeTimestampDate()
// -----------------------------------------------------------------------------
// This function will adjust the day, month, or year of the timestamp
// This change will be made regardless of packet type as it is done on the
// pcap header not the packet itself
func changeTimestampDate(packet gopacket.Packet, iDiffYear, iDiffMonth, iDiffDay int) {
	ts := packet.Metadata().CaptureInfo.Timestamp
	if iDebug == 1 {
		fmt.Println("DEBUG: Current timestamp", ts)
	}

	tsNew := ts.AddDate(iDiffYear, iDiffMonth, iDiffDay)
	if iDebug == 1 {
		fmt.Println("DEBUG: Updated timestamp", tsNew)
	}
	packet.Metadata().CaptureInfo.Timestamp = tsNew

} //changeTimestampDate()

//
// -----------------------------------------------------------------------------
// computeNeededPacketDateChange()
// -----------------------------------------------------------------------------
// Figure out if there is a change needed for the date of each packet.  We will
// compute the difference between what is in the first packet and what was passed
// in via the command line arguments.
func computeNeededPacketDateChange(year, month, day int, pcapStartTimestamp time.Time) (iDiffYear, iDiffMonth, iDiffDay int) {
	iDiffYear = 0
	iDiffMonth = 0
	iDiffDay = 0

	if year != 0 {
		iDiffYear = year - pcapStartTimestamp.Year()
	}
	if month != 0 {
		iDiffMonth = month - int(pcapStartTimestamp.Month())
	}
	if day != 0 {
		iDiffDay = day - pcapStartTimestamp.Day()
	}

	if iDebug == 1 {
		fmt.Println("DEBUG: Y/M/D deltas", iDiffYear, iDiffMonth, iDiffDay)
	}
	return
} // computeNeededPacketDateChange()

//
// -----------------------------------------------------------------------------
// changeTimestampTimeOfDay()
// -----------------------------------------------------------------------------
// This function will adjust the time of day of the timestamp
// This change will be made regardless of packet type as it is done on the
// pcap header not the packet itself
func changeTimestampTimeOfDay(packet gopacket.Packet, timeDiff time.Duration) {
	ts := packet.Metadata().CaptureInfo.Timestamp
	if iDebug == 1 {
		fmt.Println("DEBUG: Current timestamp", ts)
	}

	tsNew := ts.Add(timeDiff)
	if iDebug == 1 {
		fmt.Println("DEBUG: Updated timestamp", tsNew)
	}
	packet.Metadata().CaptureInfo.Timestamp = tsNew
} // changeTimestampTimeOfDay()

//
// -----------------------------------------------------------------------------
// computeNeededPacketTimeChange()
// -----------------------------------------------------------------------------
// Figure out if there is a change needed for the time of day of each packet.  We
// will compute the difference between what is in the first packet and what was
// passed in via the command line arguments to account for PCAP files that span
// multiple days, months, or years
func computeNeededPacketTimeChange(sTime string, pcapStartTimestamp time.Time) time.Duration {
	var amountOfTimeChange time.Duration
	var err error

	if sTime != "" {
		amountOfTimeChange, err = time.ParseDuration(sTime)
	} else {
		amountOfTimeChange, err = time.ParseDuration("0")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if iDebug == 1 {
		fmt.Println("DEBUG: Time delta", amountOfTimeChange.String())
	}
	return amountOfTimeChange
} // computeNeededPacketDateChange()

//
// -----------------------------------------------------------------------------
// getFirstPacketTimestamp
// -----------------------------------------------------------------------------
// We need to open the pcap file and read the timestamp from the first packet so
// that we can figure out an offset for all future packets.  This will address the
// problem of the pcap spanning multiple days, months, years  as we will always
// add the same amount of offset to each packet.
func getFirstPacketTimestamp(sFilename string) time.Time {
	handle, err := pcap.OpenOffline(sFilename)
	defer handle.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	_, packetHeaderInfo, _ := handle.ReadPacketData()
	ts := packetHeaderInfo.Timestamp
	if iDebug == 1 {
		fmt.Println("DEBUG: Timestamp of first packet", ts)
	}
	return ts
} // getFirstPacketTimestamp()

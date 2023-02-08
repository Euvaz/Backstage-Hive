package pkg

import (
    "fmt"
    "net"
    "net/netip"
    "strings"

    "github.com/Euvaz/Backstage-Hive/logger"
)

func parseHost(host string) (string, string) {
    // Check if host is an IP address
    address, err := netip.ParseAddr(host)
    if err != nil {
        // Attempt to resolve
        addressSlice, err := net.LookupIP(host)
        if err != nil {
            logger.Fatal(err.Error())
        }
        // Convert []net.IP to String
        addresses := fmt.Sprintf("%s", addressSlice)
        addressArray := strings.Fields(addresses[1:len(addresses)-1])
        // Ensure only a single IP was resolved
        if len(addressArray) > 1 {
            logger.Warn("More than one IP resolved; Defaulting to first address")
        }
        return addressArray[0], host
    }
    hostname := ""
    return address.String(), hostname
}

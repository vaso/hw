package service

import (
	"log"
	"net/netip"
	"slices"

	"rate_limiter/config"
)

type IPService struct {
	Blacklist []netip.Prefix
	Whitelist []netip.Prefix
}

func NewIPService(config config.AppConfig) (*IPService, error) {
	blacklist := make([]netip.Prefix, 0)
	whitelist := make([]netip.Prefix, 0)
	for _, ipMask := range config.IPBlacklist {
		network, err := netip.ParsePrefix(ipMask)
		if err != nil {
			return nil, err
		}
		blacklist = append(blacklist, network)
	}
	for _, ipMask := range config.IPWhitelist {
		network, err := netip.ParsePrefix(ipMask)
		if err != nil {
			return nil, err
		}
		whitelist = append(whitelist, network)
	}
	return &IPService{
		Blacklist: blacklist,
		Whitelist: whitelist,
	}, nil
}

func (s *IPService) AddToBlacklist(networkStr string) (bool, error) {
	network, err := netip.ParsePrefix(networkStr)
	if err != nil {
		return false, err
	}
	s.Blacklist = append(s.Blacklist, network)
	return true, nil
}

func (s *IPService) RemoveFromBlacklist(networkStr string) (bool, error) {
	network, err := netip.ParsePrefix(networkStr)
	if err != nil {
		return false, err
	}

	idx := slices.IndexFunc(s.Blacklist, func(addr netip.Prefix) bool { return addr == network })

	updatedList := make([]netip.Prefix, 0)
	updatedList = append(updatedList, s.Blacklist[:idx]...)
	updatedList = append(updatedList, s.Blacklist[idx+1:]...)
	s.Blacklist = updatedList

	return true, nil
}

func (s *IPService) AddToWhitelist(networkStr string) (bool, error) {
	network, err := netip.ParsePrefix(networkStr)
	if err != nil {
		return false, err
	}
	s.Whitelist = append(s.Whitelist, network)
	return true, nil
}

func (s *IPService) RemoveFromWhitelist(networkStr string) (bool, error) {
	network, err := netip.ParsePrefix(networkStr)
	if err != nil {
		return false, err
	}

	idx := slices.IndexFunc(s.Whitelist, func(addr netip.Prefix) bool { return addr == network })

	updatedList := make([]netip.Prefix, 0)
	updatedList = append(updatedList, s.Whitelist[:idx]...)
	updatedList = append(updatedList, s.Whitelist[idx+1:]...)
	s.Whitelist = updatedList

	return true, nil
}

func (s *IPService) IsInBlacklist(ipStr string) (bool, error) {
	log.Printf("check Blacklist: %+v %s\n", s.Blacklist, ipStr)
	return s.isInList(s.Blacklist, ipStr)
}

func (s *IPService) IsInWhitelist(ipStr string) (bool, error) {
	log.Printf("check Whitelist: %+v %s\n", s.Whitelist, ipStr)
	return s.isInList(s.Whitelist, ipStr)
}

func (s *IPService) isInList(list []netip.Prefix, ipStr string) (bool, error) {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		panic(err)
	}
	for _, network := range list {
		if network.Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

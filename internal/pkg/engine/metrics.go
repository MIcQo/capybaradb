package engine

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var createDatabaseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace:   "capybaradb",
	Subsystem:   "engine",
	Name:        "create_database_statement",
	Help:        "Number of create database statements from startup",
	ConstLabels: nil,
}, []string{})

var useCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace:   "capybaradb",
	Subsystem:   "engine",
	Name:        "use_statement",
	Help:        "Number of use statements from startup",
	ConstLabels: nil,
}, []string{})

var showCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace:   "capybaradb",
	Subsystem:   "engine",
	Name:        "show_statement",
	Help:        "Number of show statements from startup",
	ConstLabels: nil,
}, []string{})

var selectCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace:   "capybaradb",
	Subsystem:   "engine",
	Name:        "select_statement",
	Help:        "Number of select statements from startup",
	ConstLabels: nil,
}, []string{})

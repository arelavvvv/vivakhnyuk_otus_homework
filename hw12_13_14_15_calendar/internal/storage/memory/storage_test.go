package memorystorage

import (
	"testing"

	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	memstorage := New()
	err := memstorage.NewEvent(
		"title",
		"2025-02-10 13:37:00",
		"2025-02-10 18:00:00",
		"desc",
		42,
		500,
	)
	require.Nil(t, err)
	require.Len(t, memstorage.events, 1)

	err = memstorage.NewEvent(
		"title",
		"2025-02-10 14:37:00",
		"2025-02-10 18:00:00",
		"desc",
		42,
		500,
	)

	require.Error(t, err)
	require.ErrorIs(t, err, storage.DateBusyError{
		Title:    "title",
		DateFrom: "2025-02-10 13:37:00",
		DateTo:   "2025-02-10 18:00:00",
	},
	)

	events, err := memstorage.GetEvents()
	require.Nil(t, err)
	require.Len(t, events, 1)

	err = memstorage.UpdEvent(
		events[0].ID,
		"updated",
		events[0].DateTimeStart,
		events[0].DateTimeEnd,
		events[0].Description,
		events[0].UserID,
		events[0].Delay,
	)
	require.Nil(t, err)
	require.Equal(t, memstorage.events[0].Title, "updated")

	err = memstorage.EventDelete(events[0].ID)
	require.Nil(t, err)

	events, err = memstorage.GetEvents()
	require.Nil(t, err)
	require.Len(t, events, 0)
}

func TestStorageListing(t *testing.T) {
	memstorage := New()

	err := memstorage.NewEvent(
		"title",
		"2025-02-10 13:37:00",
		"2025-02-10 18:00:00",
		"description",
		42,
		500,
	)
	require.Nil(t, err)

	err = memstorage.NewEvent(
		"title",
		"2025-02-11 13:37:00",
		"2025-02-11 18:00:00",
		"description",
		42,
		500,
	)
	require.Nil(t, err)

	events, _ := memstorage.GetEventsDay("2025-02-10")
	require.Len(t, events, 1)

	events, _ = memstorage.GetEventsWeek("2025-02-10")
	require.Len(t, events, 2)

	events, _ = memstorage.GetEventsMonth("2025-02-10")
	require.Len(t, events, 2)
}

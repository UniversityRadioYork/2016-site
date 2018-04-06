$(document).ready(function() {
    timeslotTable = $('#timeslots').DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
        "columnDefs": [
        {
            "targets": [ 0 ],
            "visible": false,
            "searchable": true
        },
        {
            "targets": [ 1,2 ],
            "searchable": false
        },

        ]
    });
    $('.timeslotFilterSeason').on( 'click', 'a', function () {
        $('#timeslotLatest').hide();

        timeslotTable
            .columns( 0 )
            .search( $(this).data("seasonid") )
            .draw();

        $('#timeslotList').show();
    } );
} );
